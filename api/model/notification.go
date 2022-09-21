package model

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	webpush "github.com/SherClockHolmes/webpush-go"
)

//NotificationType represents a notification type
type NotificationType string

const (
	SMSNotificationType   NotificationType = "sms"
	EmailNotificationType NotificationType = "email"
)

//Notification represents a general notification
type Notification struct {
	Owner            primitive.ObjectID `bson:"owner_id"  json:"owner_id" groups:"private"`
	Message          string             `bson:"message"  json:"message" groups:"private"`
	Type             NotificationType   `bson:"type"  json:"type" groups:"public"`
	Created          time.Time          `bson:"created_at"  json:"created_at" groups:"public"`
	SentAt           time.Time          `bson:"sent_at"  json:"sent_at" groups:"public"`
	Delivered        bool               `bson:"delivered"  json:"delivered" groups:"public"`
	Log              string             `bson:"log"  json:"log" groups:"protected"`
	Gateway          string             `bson:"gateway"  json:"gateway" groups:"private"`
	Priority         bool               `bson:"priority"  json:"priority" groups:"protected"`
	IsApplicationLog bool               `bson:"app_log"  json:"app_log" groups:"protected"`
}

type WebNotification struct {
	Mobile    webpush.Subscription `bson:"mobile" json:"-" groups:"private"`
	Secondary webpush.Subscription `bson:"secondary" json:"-" groups:"private"`
}

//SMSNotification models a general notification SMS
//See SMSOTP for auth relate notification
type SMSNotification struct {
	Notification
	PhoneNumber string `bson:"phone_number"  json:"phone_number" groups:"private"`
}

//EmailNotification models a general notification email
//unlike SMSNotification, each EmailNotification has a call-to-action
type EmailNotification struct {
	Notification
	Email      string            `bson:"email" json:"email" groups:"private"`
	Title      string            `bson:"title"  json:"title" groups:"private"`
	Template   string            `bson:"template"  json:"template" groups:"private"`
	ActionURL  string            `bson:"action_url"  json:"action_url" groups:"private"`
	ActionName string            `bson:"action_name"  json:"action_name" groups:"private"`
	Dictionary map[string]string `bson:"dictionary" json:"dictionary" groups:"private"`
}

func (user *User) SetPushSubscription(sub []byte, secondary bool) (err error) {
	s := &webpush.Subscription{}
	json.Unmarshal(sub, s)

	n := user.WebNotification
	if secondary {
		n.Secondary = *s
	} else {
		n.Mobile = *s
		n.Secondary = *s
	}

	log.Println("we got ", s)

	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"web_notifications": n,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	log.Println(err)
	return
}

func AddProcessedNotification(notification interface{}) (err error) {
	_, err = notificationsCollection().InsertOne(context.TODO(), notification)
	return
}
