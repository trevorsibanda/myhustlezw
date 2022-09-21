package cron

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
)

var (
	cron       *gocron.Scheduler
	rabbitConn *amqp.Connection
)

type thumbnailDimensions []int

var (
	paymentsSettlerInterval = os.Getenv("CRON_PAYMENTS_SETTLER_INTERVAL")
	thumbnailRoleFormats    = map[string][]thumbnailDimensions{
		"profile_pic": {
			{120, 120},
			{240, 240},
			{480, 480},
		},
		"content": {
			{240, 240},
			{120, 120},
			{480, 480},
			{600, 600},
		},
		"headline": {
			{600, 600},
			{480, 480},
			{120, 120},
		},
		"cover": {
			{480, 480},
			{930, 389},
			{720, 195},
			{1080, 300},
			{740, 145},
		},
		"feature": {
			{480, 480},
			{600, 600},
			{720, 720},
		},
	}
)

//INIT RABBITMQ CONN
func InitializeMessagingQueue() {
	var err error
	url := os.Getenv("RABBITMQ_DIAL_URL")
	log.Printf("messaging: Starting up with connection: %s", url)
	if rabbitConn, err = amqp.Dial(url); err != nil {
		panic(fmt.Sprintf("messaging: Failed to connect to RabbitMQ: %v", err))
	} else {
		if _, err := rabbitConn.Channel(); err != nil {
			panic(fmt.Sprintf("messaging: Failed to open a channel: %v", err))
		}

	}

	//rabbitConn.NotifyClose()

	//declareRabbitQueue(_emailsPendingQueue, emailsPendingChannel, emailsPendingQueue)
	//declareRabbitQueue(_emailsPriorityQueue, emailsPriorityChannel, emailsPriorityQueue)
	emailsPendingChannel, emailsPendingQueue = declareRabbitQueue(_emailsPendingQueue)
	emailsPriorityChannel, emailsPriorityQueue = declareRabbitQueue(_emailsPriorityQueue)

	audioProcessingChannel, audioProcessingQueue = declareRabbitQueue(_audioProcessingQueue)
	videoProcessingChannel, videoProcessingQueue = declareRabbitQueue(_videoProcessingQueue)
	thumbnailGenChannel, thumbnailGenQueue = declareRabbitQueue(_thumbnailGenQueue)
	paymentsPendingChannel, paymentsPendingQueue = declareRabbitQueue(_paymentsPendingQueue)

	log.Printf("messaging: Connected to RabbitMQ")

}

func declareRabbitQueue(queueName string) (ch *amqp.Channel, q *amqp.Queue) {
	var err error
	ch, err = rabbitConn.Channel()
	if err != nil {
		panic(fmt.Sprintf("messaging: Failed to open a channel: %v", err))
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(fmt.Sprintf("emails Failed to set QoS %v", err))
	}

	if qq, err1 := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		panic(fmt.Sprintf("messaging: Failed to declare a queue: %v", err1))
	} else {
		q = &qq
	}
	assert(q != nil, "Queue is nil")
	assert(ch != nil, "Channel is nil")

	return
}

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

//Init initializes the cron service
func Init() {

	cronNode := os.Getenv("MYHUSTLE_CRON_NODE")
	InitializeMessagingQueue()

	cron = gocron.NewScheduler(time.UTC)

	if cronNode == "true" || cronNode == "1" {
		log.Printf("Starting cron service on node %s", cronNode)
		log.Printf("CRON: PaymentSettler running every %s", paymentsSettlerInterval)
		cron.Every(paymentsSettlerInterval).Do(PaymentsSettlerRunner)
	} else {
		log.Printf("Cron disabled on node %s", cronNode)
	}

	cron.Every("30m").Do(UpdateExchangeRate)
	UpdateExchangeRate()

	InitializeNotificationClients()
	cron.StartAsync()
	go videoProcessRunner()
	go thumbnailGeneratorRunner()
	go paymentsDispatcherRunner()
	go audioProcessRunner()
	go emailsDispatcherRunner()
	go priorityEmailsRunner()
}
