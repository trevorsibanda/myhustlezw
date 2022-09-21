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
	_audioProcessingQueue  = "audio-processing-queue"
	audioProcessingQueue   *amqp.Queue
	audioProcessingChannel *amqp.Channel
)

type audioProcessItem struct {
	Creator model.User
	Node    string
	Audio   model.CreatorFile
}

//ScheduleaudioEncoding adds a audio to an encoding queue
func ScheduleAudioEncoding(creator model.User, audio model.CreatorFile) {
	log.Printf("audioprocessing Added  %v by @%v to queue ", audio, creator.Username)
	entry := audioProcessItem{
		Creator: creator,
		Node:    runnerNode,
		Audio:   audio,
	}
	if jsonEntry, err := json.Marshal(entry); err != nil {
		log.Printf("audioprocessing Failed to marshal entry %v", err)
	} else {
		if err := audioProcessingChannel.Publish("", audioProcessingQueue.Name, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonEntry,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			log.Printf("audioprocessing Failed to publish to queue %v", err)
			//send support email
			go processAudio(entry)

			//todo: optionally start processing the video
			SendEngineerSupportEmail(fmt.Sprintf("Failed to publish to queue %v", err), entry)
		}
	}
	log.Printf("audioprocessing Added  %v by @%v to queue ", audio, creator.Username)
}

//audioProcessRunner dequeues an item from the audioQueue and processes it
func audioProcessRunner() {
	log.Printf("audioprocessing Started at %s", time.Now())
	assert(audioProcessingQueue != nil, "audioProcessing queue is nil")

	defer audioProcessingChannel.Close()

	msgs, err := audioProcessingChannel.Consume(
		_audioProcessingQueue, // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		panic(fmt.Sprintf("audioprocessing Failed to consume from queue %v", err))
	}
	for d := range msgs {

		var item audioProcessItem
		if err := json.Unmarshal([]byte(d.Body), &item); err != nil {
			log.Printf("audioprocessing Failed to unmarshal message %v", err)
			break
		}
		log.Printf("audioprocessing Processing item %v", item.Audio.ID)
		processAudio(item)
		d.Ack(false)
		log.Printf("audioprocessing Waiting for item to process")

	}
}

func processAudio(item audioProcessItem) {
	log.Printf("audioprocessing Processing  %v by @%v to queue ", item.Audio, item.Creator.Username)

	if hlsFile, err := storage.ConvertAudioToAdaptiveHLS(item.Audio); err == nil {
		file := item.Audio
		file.ETag = time.Now().String()
		file.Filename = hlsFile
		file.Mime = "application/vnd.apple.mpegurl"
		file.Processed = true
		//save to db
		err = file.SaveUploadToStorageChanges()
		log.Println(err)
	} else {
		log.Printf("audioprocessing Processed audio. Return errors: %v ", err)
	}
}
