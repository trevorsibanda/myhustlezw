package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/streadway/amqp"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/storage"
)

var (
	_thumbnailGenQueue                = "image:thubnail-gen"
	thumbnailGenQueue   *amqp.Queue   = nil
	thumbnailGenChannel *amqp.Channel = nil
	runnerNode                        = os.Getenv("MYHUSTLE_NODE_ID")
)

type thumbnailProcessItem struct {
	Creator model.User
	Node    string
	Image   model.CreatorFile
}

//ScheduleGenerateThumbnails schedules  thumbnail generation
func ScheduleGenerateThumbnails(creator model.User, image model.CreatorFile) {
	entry := thumbnailProcessItem{
		Creator: creator,
		Node:    runnerNode,
		Image:   image,
	}

	log.Printf("thumbnailgen Added  %v by @%v to queue ", image, creator.Username)
	if jsonEntry, err := json.Marshal(entry); err != nil {
		log.Printf("thumbnailGen Failed to marshal entry %v", err)
	} else {
		if err := thumbnailGenChannel.Publish("", thumbnailGenQueue.Name, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonEntry,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			log.Printf("thumbnailGen Failed to publish to queue %v", err)
			//send support email

			go processThumbnail(entry)

			//todo: optionally start processing the video
			SendEngineerSupportEmail(fmt.Sprintf("Failed to publish to queue %v", err), entry)
		}
	}

}

//thumbnail generator runner generates thumbnails for given image in queue
func thumbnailGeneratorRunner() {

	log.Printf("thumbnailgen Started at %s", time.Now())
	assert(thumbnailGenQueue != nil, "thumbnailgen queue is nil")

	defer thumbnailGenChannel.Close()

	msgs, err := thumbnailGenChannel.Consume(
		_thumbnailGenQueue, // queue
		"",                 // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		panic(fmt.Sprintf("thumbnailGen Failed to consume from queue %v", err))
	}
	for d := range msgs {

		var item thumbnailProcessItem

		if err := json.Unmarshal([]byte(d.Body), &item); err != nil {
			log.Printf("thumbnailgen Failed to unmarshal message %v", err)
			break
		}
		log.Printf("thumbnailgen Processing item %v", item.Image.ID)
		processThumbnail(item)
		d.Ack(false)
		log.Printf("thumbnailgen Waiting for item to process")

	}
}

func processThumbnail(item thumbnailProcessItem) {
	src := storage.LocalStoragePath(item.Image.Filename)

	var rDimensions []thumbnailDimensions
	var ok bool

	image, err := imaging.Open(src)

	if err != nil {
		log.Printf("thumbnailgen Failed to load image to process %v Reason: %v ", item.Image, err)
		return
	}

	if rDimensions, ok = thumbnailRoleFormats[item.Image.Role]; !ok {
		rDimensions = thumbnailRoleFormats["content"]
	}

	for _, res := range rDimensions {
		resized := imaging.Resize(image, res[0], res[1], imaging.Lanczos)
		if resized != nil {
			blurred := imaging.Blur(resized, 30.8)

			target := strings.Split(item.Image.Filename, ".")
			dest := storage.LocalStoragePath(fmt.Sprintf("%s_%dx%d.png", target[0], res[0], res[1]))
			imaging.Save(resized, dest)

			if blurred != nil {
				dest := storage.LocalStoragePath(fmt.Sprintf("%s_%dx%d_blurred.png", target[0], res[0], res[1]))
				imaging.Save(blurred, dest)
			}
		} else {
			log.Printf("thumbnailgen Failed to resize image %v %v %v", item.Image, image, resized)
		}
	}
	log.Printf("thumbnailgen Finished processing at %v", time.Now())
}
