package model

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VisibilityType string

const (
	PublicVisible  VisibilityType = "public"
	PrivateVisible VisibilityType = "private"
	PaidVisible    VisibilityType = "paid"
)

type SubscriptionPeriod string

const (
	PeriodMonth   SubscriptionPeriod = "month"
	PeriodDay     SubscriptionPeriod = "day"
	PeriodWeek    SubscriptionPeriod = "week"
	PeriodYear    SubscriptionPeriod = "year"
	PeriodForever SubscriptionPeriod = "forever"
)

//GetSubscriptions returns all of a fans active subscriptions
func (fan User) GetSubscriptions(filterString string) (subs []CreatorSupport, err error) {
	subs = make([]CreatorSupport, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"fan_id":       fan.ID,
		"support_type": SubscribeToCreatorAccount,
	}

	switch strings.ToLower(filterString) {
	case "active":
		filter["expires"] = bson.M{"$gt": time.Now()}
	case "expired":
		filter["expires"] = bson.M{"$lte": time.Now()}
	}

	opts := options.FindOptions{}
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err = supportersCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &subs)

	return
}

//GetSubscription gets a subscription to a particular creatorID
func (fan User) GetSubscription(creatorID primitive.ObjectID) (sub CreatorSupport, err error) {
	filter := bson.M{
		"fan_id":       fan.ID,
		"creator_id":   creatorID,
		"support_type": SubscribeToCreatorAccount,
	}

	opts := options.FindOneOptions{}
	opts.SetSort(bson.M{"created_at": -1})

	err = supportersCollection().FindOne(context.TODO(), filter, &opts).Decode(&sub)
	return
}
