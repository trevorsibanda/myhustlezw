package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/storage"
)

var (
	_videoProcessingQueue  = "video-processing-queue"
	videoProcessingQueue   *amqp.Queue
	videoProcessingChannel *amqp.Channel
)

type videoProcessItem struct {
	Creator model.User
	Node    string
	Video   model.CreatorFile
}

//ScheduleVideoEncoding adds a video to an encoding queue
func ScheduleVideoEncoding(creator model.User, video model.CreatorFile) {
	log.Printf("videoprocessing Added  %v by @%v to queue ", video, creator.Username)
	entry := videoProcessItem{
		Creator: creator,
		Node:    runnerNode,
		Video:   video,
	}
	if jsonEntry, err := json.Marshal(entry); err != nil {
		log.Printf("videoprocessing Failed to marshal entry %v", err)
	} else {
		if err := videoProcessingChannel.Publish("", videoProcessingQueue.Name, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonEntry,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			log.Printf("videoprocessing Failed to publish to queue %v", err)
			//send support email

			go processVideo(entry)

			//todo: optionally start processing the video
			SendEngineerSupportEmail(fmt.Sprintf("Failed to publish to queue %v", err), entry)
		}
	}
	log.Printf("videoprocessing Added  %v by @%v to queue ", video, creator.Username)
}

//VideoProcessRunner dequeues an item from the videoQueue and processes it
func videoProcessRunner() {
	log.Printf("videoprocessing Started at %s", time.Now())
	assert(videoProcessingQueue != nil, "videoProcessing queue is nil")

	defer videoProcessingChannel.Close()

	msgs, err := videoProcessingChannel.Consume(
		_videoProcessingQueue, // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		panic(fmt.Sprintf("videoprocessing Failed to consume from queue %v", err))
	}
	for d := range msgs {

		var item videoProcessItem
		if err := json.Unmarshal([]byte(d.Body), &item); err != nil {
			log.Printf("videoprocessing Failed to unmarshal message %v", err)
			break
		}
		log.Printf("videoprocessing Processing item %v", item.Video.ID)
		processVideo(item)
		d.Ack(false)
		log.Printf("videoprocessing Waiting for item to process")

	}
}

func processVideo(item videoProcessItem) {

	log.Printf("videoprocessing Processing  %v by @%v to queue ", item.Video, item.Creator.Username)

	if hlsFile, err := storage.ConvertVideoToAdaptiveHLS(item.Video); err == nil {
		file := item.Video
		file.ETag = time.Now().String()
		file.Filename = hlsFile
		file.Mime = "application/vnd.apple.mpegurl"
		file.Processed = true
		//save to db
		err = file.SaveUploadToStorageChanges()
		log.Println(err)
	} else {
		log.Printf("videoprocessing Processed video. Return errors: %v ", err)
	}

}
