package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedItem struct {
	URI          string    `json:"uri" groups:"public"`
	Username     string    `json:"username" groups:"public"`
	CreatedAt    time.Time `json:"created_at" groups:"public"`
	Thumbnail    string    `json:"thumbnail" groups:"public"`
	PreviewURL   string    `json:"preview_url" groups:"public"`
	CanView      bool      `json:"can_view" groups:"public"`
	Type         string    `json:"type" groups:"public"`
	Subscription string    `json:"subscription" groups:"public"`
	Price        float64   `json:"price" groups:"public"`
	Title        string    `json:"title" groups:"public"`
	Description  string    `json:"description" groups:"public"`
}

type RawFeedItem struct {
	Creator User     `json:"creator" groups:"public"`
	Content Campaign `json:"content" groups:"public"`
}

func NewFeedItem(rf RawFeedItem, sessionKey *string, canView bool) FeedItem {
	creator := rf.Creator
	campaign := rf.Content
	return FeedItem{
		URI:          fmt.Sprintf("/@%s/%s", creator.Username, campaign.URI),
		Username:     creator.Username,
		CreatedAt:    campaign.Created,
		Thumbnail:    campaign.PreviewURL(sessionKey, canView, 240, 240),
		PreviewURL:   campaign.PreviewURL(sessionKey, canView, 480, 480),
		CanView:      canView,
		Type:         campaign.Type,
		Subscription: campaign.Subscription,
		Price:        campaign.Price,
		Title:        campaign.Title,
		Description:  campaign.Description,
	}
}

type FeedView int

func ViewFromString(str string) (v FeedView) {
	switch strings.ToLower(str) {
	case "all":
		v = FeedViewAll
	case "explore":
		v = FeedViewDiscover
	case "filtered":
		v = FeedViewFiltered
	}
	return
}

var FeedViewAll FeedView = 0
var FeedViewDiscover FeedView = 1
var FeedViewFiltered FeedView = 2

type FeedFilter int

func FilterFromString(str string) (f FeedFilter) {
	switch strings.ToLower(str) {
	case "pay_per_view":
		f = FeedFilterPaidContent
	case "subscriptions":
		f = FeedFilterSubscriptionsOnly
	case "creator":
		f = FeedFilterByCreator
	default:
		f = NoFilter
	}
	return
}

var NoFilter FeedFilter = -1
var FeedFilterPaidContent FeedFilter = 0
var FeedFilterSubscriptionsOnly FeedFilter = 1
var FeedFilterByCreator FeedFilter = 2

//TODO: replace this with gorse.io content recommendation engine when resources allow
func (user User) GenerateUserFeed(skip int64, view FeedView, filter FeedFilter, paidContentIDs, subscribedUserIDs, interactedUsers []primitive.ObjectID, creatorFilter *primitive.ObjectID) (feed []RawFeedItem, err error) {
	dbFilter := bson.M{
		"active": bson.M{"$ne": false},
	}

	dbSort := bson.M{"created_at": -1}

	switch view {
	case FeedViewAll:
		concatUsers := append(subscribedUserIDs, interactedUsers...)
		dbFilter = bson.M{"$or": bson.A{
			bson.M{"owner_id": bson.M{"$in": concatUsers}},
			bson.M{"_id": bson.M{"$in": paidContentIDs}},
		},
		}

	case FeedViewDiscover:
		dbFilter["discoverable"] = bson.M{"$ne": false}

	case FeedViewFiltered:
		switch filter {
		case FeedFilterPaidContent:
			dbFilter["_id"] = bson.M{"$in": paidContentIDs}
		case FeedFilterSubscriptionsOnly:
			dbFilter["owner_id"] = bson.M{"$in": subscribedUserIDs}
			dbFilter["subscription"] = "fans"
		case FeedFilterByCreator:
			dbFilter["owner_id"] = creatorFilter
		}
	}

	opts := options.FindOptions{}

	opts.SetSort(dbSort)
	opts.SetSkip(skip)

	var campaigns []Campaign
	cursor, err := campaignsCollection().Find(context.TODO(), dbFilter, &opts)
	if err != nil {
		return
	}
	err = cursor.All(context.TODO(), &campaigns)
	if err != nil {
		return
	}

	userIds := make(map[primitive.ObjectID]User, 0)
	for _, c := range campaigns {
		u, ok := userIds[c.Owner]
		if !ok {
			u, _ = RetrieveCreatorByID(c.Owner)
			userIds[c.Owner] = u
		}
		feed = append(feed, RawFeedItem{Creator: u, Content: c})
	}

	return
}
