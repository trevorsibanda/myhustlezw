package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreatorSupportType = string

const (
	SubscribeToCreatorAccount CreatorSupportType = "subscribed"
	SupportCreator            CreatorSupportType = "support"
	PayPerViewAccess          CreatorSupportType = "paid_content"
	ServiceRequest            CreatorSupportType = "service_request"
)

//FilledServiceForm stores information filled for a service campaigns
type FilledServiceForm struct {
	Created           time.Time `bson:"created" json:"created" groups:"public"`
	Fulfilled         bool      `bson:"fulfilled" json:"fulfilled" groups:"public"`
	FulfilledAt       time.Time `bson:"fulfilled_at" json:"fulfilled_at" groups:"public"`
	FulfillBy         time.Time `bson:"fulfill_by" json:"fulfill_by" groups:"public"`
	RefundedAt        time.Time `bson:"refunded_at" json:"refunded_at" groups:"public"`
	RefundBy          time.Time `bson:"refund_by" json:"refund_by" groups:"public"`
	Refunded          bool      `bson:"refund" json:"refunded" groups:"public"`
	Email             string    `bson:"email" json:"email" groups:"public"`
	Phone             string    `bson:"phone" json:"phone" groups:"public"`
	Fullname          string    `bson:"fullname" json:"fullname" groups:"public"`
	Question          string    `bson:"question"  json:"question" groups:"public"`
	Answer            string    `bson:"answer" json:"answer" groups:"public"`
	Instructions      string    `bson:"instructions"  json:"instructions" groups:"public"`
	OptionQuestion    string    `bson:"options_title"  json:"options_title" groups:"public"`
	OptionValues      []string  `bson:"option_values"  json:"option_values" groups:"public"`
	SelectedOption    string    `bson:"selected_option" json:"selected_option" groups:"public"`
	ThankYouMessage   string    `bson:"thanks_message"  json:"thanks_message" groups:"public"`
	ExtraInfoQuestion string    `bson:"extra_info_question"  json:"extra_info_question" groups:"public"`
	AllowExtraInfo    bool      `bson:"allow_extra_info"  json:"allow_extra_info" groups:"private"`
	ExtraInfo         string    `bson:"extra_info" json:"extra_info" groups:"public"`
	QuantityLeft      int       `bson:"quantity_left"  json:"quantity_left" groups:"private"`
	Log               string    `bson:"log" json:"log" groups:"protected"`
}

func (form FilledServiceForm) Dict() map[string]string {
	m := map[string]string{
		"Email":        form.Email,
		"Phone":        form.Phone,
		"Fullname":     form.Fullname,
		form.Question:  form.Answer,
		"Instructions": form.Instructions,
	}
	if form.Fulfilled {
		m["Order Status"] = fmt.Sprintf("Fulfilled at %s", form.FulfilledAt.Format(time.UnixDate))
	} else if form.Refunded {
		m["Order Status"] = fmt.Sprintf("Refunded at %s", form.RefundedAt.Format(time.UnixDate))
	} else {
		m["Order Status"] = "Pending"
		m["Fulfill by"] = form.FulfillBy.Format(time.UnixDate)
		m["Refund by"] = form.RefundBy.Format(time.UnixDate)
	}
	return m
}

//ListPaidCampaigns list all paid campaigns
func (user User) ListPaidCampaigns() (list []primitive.ObjectID, err error) {
	list = make([]primitive.ObjectID, 0)
	res := make([]CreatorSupport, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"fan_id":       user.ID,
		"support_type": PayPerViewAccess,
	}
	log.Println(user.ID)
	cursor, err = supportersCollection().Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &res)

	for _, p := range res {
		list = append(list, p.Campaign)
	}

	return
}

//ListRecentSubmittedForms list all submitted forms. A submitted form is also paid for.
func (campaign Campaign) ListRecentSubmittedForms(maxItems int64, skip int64) (list []FilledServiceForm, err error) {
	list = make([]FilledServiceForm, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"campaign_id": campaign.ID,
	}
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &maxItems,
	}

	cursor, err = supportersCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &list)

	return
}

//GetSupporter returns a supporter by ID
func GetSupporter(id primitive.ObjectID, creatorID *primitive.ObjectID) (supporter CreatorSupport, err error) {
	filter := bson.M{
		"_id": id,
	}

	if creatorID != nil {
		filter["creator_id"] = creatorID
	}

	err = supportersCollection().FindOne(context.TODO(), filter).Decode(&supporter)
	return
}

//Subscription models a subscription to a creator's content
type Subscription struct {
	ID            primitive.ObjectID `bson:"_id"  json:"_id" groups:"public"`
	Status        string             `bson:"status" json:"status" groups:"public"`
	CancelledWhen time.Time          `bson:"cancelled_when"  json:"cancelled_when" groups:"public"`
	CancelReason  string             `bson:"cancel_reason"  json:"cancel_reason" groups:"public"`
	Expires       time.Time          `bson:"expires_at"  json:"expires_at" groups:"public"`
}

//CreatorSupport models an interaction with a creator's content which can be considered supporting the content
//This includes leaving a tip, subscribing or paying for a service
type CreatorSupport struct {
	ID          primitive.ObjectID `bson:"_id"  json:"_id" groups:"public"`
	UnlockCode  primitive.ObjectID `bson:"unlock_code" json:"unlock_code,omitempty" groups:"protected"`
	Type        CreatorSupportType `bson:"support_type"  json:"support_type" groups:"public"`
	Creator     primitive.ObjectID `bson:"creator_id"  json:"creator_id" groups:"public"`
	Campaign    primitive.ObjectID `bson:"campaign_id"  json:"campaign_id" groups:"private"`
	Supporter   primitive.ObjectID `bson:"fan_id"  json:"fan_id" groups:"public"`
	WalletOp    primitive.ObjectID `bson:"wallet_operation_id"  json:"wallet_operation_id" groups:"private"`
	Created     time.Time          `bson:"created_at"  json:"created_at" groups:"public"`
	Expires     time.Time          `bson:"expires" json:"expires" groups:"public"`
	LastUpdated time.Time          `bson:"last_updated"  json:"last_updated" groups:"private"`
	DisplayName string             `bson:"display_name"  json:"display_name" groups:"public"`
	Email       string             `bson:"email" json:"email" groups:"protected"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number" groups:"protected"`
	IsAnonymous bool               `bson:"is_anonymous"  json:"is_anonymous" groups:"public"`
	Currency    string             `bson:"currency" json:"currency" groups:"private"`
	Items       int                `bson:"items" json:"items" groups:"public"`
	ItemName    string             `bson:"item_name,omitempty" json:"item_name,omitempty" groups:"public"`
	Amount      float64            `bson:"amount"  json:"amount" groups:"private"`
	Comment     string             `bson:"comment" json:"comment" groups:"public"`
	Hidden      bool               `bson:"hidden" json:"hidden" groups:"private"`
	HideMessage bool               `bson:"hide_message" json:"hide_message" groups:"private"`
	Form        *FilledServiceForm `bson:"form,omitempty" json:"form,omitempty" groups:"private"`
}

func (s CreatorSupport) DashboardURL() string {
	return fmt.Sprintf("%s://%s/creator/supporters/%s", myhustleScheme, myhustleEndpoint, s.ID.Hex())
}

//SaveChanges saves the uploaded files and updates the changes
func (supporter CreatorSupport) SaveChanges() (err error) {
	filter := bson.M{
		"_id": supporter.ID,
	}

	_, err = supportersCollection().ReplaceOne(context.TODO(), filter, supporter)
	return
}

func (s CreatorSupport) ContentURL() string {
	return fmt.Sprintf("%s://%s/_unlock/%s", myhustleScheme, myhustleEndpoint, s.UnlockCode.Hex())
}

//ListRecentSupporters list all supporters for a creator
func (creator User) ListRecentSupporters(maxItems int64, skip int64, campaign *primitive.ObjectID, showHiddenMessages bool) (list []CreatorSupport, err error) {
	list = make([]CreatorSupport, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"creator_id": creator.ID,
		"hidden":     bson.M{"$ne": true},
	}
	if showHiddenMessages {
		delete(filter, "hidden")
	}
	if campaign != nil {
		filter["campaign_id"] = *campaign
	}

	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &maxItems,
	}

	opts.SetSort(bson.M{"created_at": -1})

	cursor, err = supportersCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &list)

	return
}

//ListRecentSupporters list all supporters for a campaign
func (campaign Campaign) ListRecentSupporters(maxItems int64, skip int64) (list []CreatorSupport, err error) {
	list = make([]CreatorSupport, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"$and": []bson.M{
			{"creator_id": campaign.Owner},
			{"campaign_id": campaign.ID},
			{"active": true},
		},
	}
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &maxItems,
	}

	cursor, err = supportersCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &list)

	return
}

//AddSubscriber creates a new subscriber for a given supporter
func (creator User) AddSubscriber(supporter User, walletOperation WalletOperation, campaignID primitive.ObjectID, shareDetails bool) (s CreatorSupport, err error) {
	s = CreatorSupport{
		ID:          primitive.NewObjectID(),
		Email:       supporter.Email,
		PhoneNumber: supporter.PhoneNumber,
		Creator:     creator.ID,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Expires:     time.Now().Add(time.Hour * 24 * 31 * 3), //subscribe for 3 months
		Type:        SubscribeToCreatorAccount,
		Supporter:   supporter.ID,
		Campaign:    campaignID,
		DisplayName: supporter.Fullname,
		Items:       1,
		Currency:    walletOperation.Currency,
		ItemName:    creator.Username,
		IsAnonymous: !shareDetails,
		Amount:      walletOperation.Amount,
		WalletOp:    walletOperation.ID,
	}

	_, err = supportersCollection().InsertOne(context.TODO(), s)
	return
}

//AddSupportMessage adds a new support message
func (creator User) AddSupportMessage(supporter User, message string, walletOperation WalletOperation, items int, itemName string, shareDetails bool) (s CreatorSupport, err error) {
	dname := "Anonymous"
	if shareDetails && supporter.Fullname != "" {
		dname = supporter.Fullname
	}
	s = CreatorSupport{
		ID:          primitive.NewObjectID(),
		Email:       supporter.Email,
		PhoneNumber: supporter.PhoneNumber,
		Type:        SupportCreator,
		Creator:     creator.ID,
		Supporter:   supporter.ID,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Expires:     time.Now().Add(time.Minute * 60 * 24 * 365), //one year subscription. change this to be dynamic
		DisplayName: dname,
		IsAnonymous: !shareDetails,
		Amount:      walletOperation.Amount,
		Items:       items,
		ItemName:    itemName,
		Currency:    walletOperation.Currency,
		WalletOp:    walletOperation.ID,
		Comment:     message,
	}

	_, err = supportersCollection().InsertOne(context.TODO(), s)
	return
}

//AddSupportMessage adds a new support message
func (creator User) AddPayPerViewAccess(campaign Campaign, supporter User, walletOperation WalletOperation) (s CreatorSupport, err error) {
	s = CreatorSupport{
		ID:          primitive.NewObjectID(),
		UnlockCode:  primitive.NewObjectID(),
		Email:       supporter.Email,
		PhoneNumber: supporter.PhoneNumber,
		Type:        PayPerViewAccess,
		Creator:     creator.ID,
		Campaign:    campaign.ID,
		Items:       1,
		ItemName:    campaign.Title,
		Supporter:   supporter.ID,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Expires:     time.Now().Add(time.Minute * 60 * 24 * 365), //one year subscription. change this to be dynamic
		DisplayName: supporter.Fullname,
		IsAnonymous: supporter.ID.IsZero(),
		Amount:      walletOperation.Amount,
		Currency:    walletOperation.Currency,
		WalletOp:    walletOperation.ID,
	}

	_, err = supportersCollection().InsertOne(context.TODO(), s)
	return
}

//AddServiceRequest adds a request for a service
func (creator User) AddServiceRequest(supporter User, service Campaign, filledForm *FilledServiceForm, payment PendingPayment, walletOp WalletOperation) (s CreatorSupport, err error) {
	dname := "Someone"
	if !payment.Anonymous {
		dname = supporter.Fullname
	}
	s = CreatorSupport{
		ID:          primitive.NewObjectID(),
		UnlockCode:  primitive.NewObjectID(),
		Email:       supporter.Email,
		PhoneNumber: supporter.PhoneNumber,
		Type:        ServiceRequest,
		Creator:     creator.ID,
		Campaign:    service.ID,
		Supporter:   supporter.ID,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		DisplayName: dname,
		IsAnonymous: payment.Anonymous,
		Amount:      walletOp.Amount,
		Items:       1,
		ItemName:    service.Title,
		Currency:    walletOp.Currency,
		WalletOp:    walletOp.ID,
		Form:        filledForm,
	}
	_, err = supportersCollection().InsertOne(context.TODO(), s)
	return
}

//FetchSupporter fetches a supporter for a given user
func (campaign Campaign) FetchSupporters(user primitive.ObjectID, sType CreatorSupportType) (list []CreatorSupport, err error) {
	list = make([]CreatorSupport, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"$and": []bson.M{
			{"creator_id": campaign.Owner},
			{"campaign_id": campaign.ID},
			{"type": sType},
			{"active": true},
		},
	}

	cursor, err = supportersCollection().Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &list)

	return
}
