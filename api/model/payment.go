package model

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentGateway string
type PaymentCurrency string
type PendingPaymentAction string

const (
	ZWL PaymentCurrency = "ZWL"
	USD PaymentCurrency = "USD"
)

const (
	SubscribeOnPaid           PendingPaymentAction = "grant_subscribe"
	AllowAccessOnPaid         PendingPaymentAction = "grant_campaign"
	AllowServiceRequestOnPaid PendingPaymentAction = "grant_service"
	AddSupportMessageOnPaid   PendingPaymentAction = "add_supporter_message"
	VerifyAccountOnPaid       PendingPaymentAction = "verify_account"
)

const (
	GatewayPaynow    PaymentGateway = "paynow"
	GatewayEcocash   PaymentGateway = "ecocash"
	Gateway2Checkout PaymentGateway = "2checkout"
)

type PaymentLogInitiator = string

const (
	PaymentInitiatedServer  PaymentLogInitiator = "server"
	PaymentInitiatedGateway PaymentLogInitiator = "gateway"
)

type PaymentLog struct {
	Gateway        PaymentGateway      `bson:"gateway"`
	Request        interface{}         `bson:"request"`
	Response       interface{}         `bson:"response"`
	Initiator      PaymentLogInitiator `bson:"initiator"`
	TargetCreator  primitive.ObjectID  `bson:"creator"`
	TargetFan      primitive.ObjectID  `bson:"fan"`
	TargetCampaign primitive.ObjectID  `bson:"campaign"`
	WalletOp       primitive.ObjectID  `bson:"wallet_op"`
}

//LogPayment logs a payment
func LogPayment(entry PaymentLog) (err error) {
	_, err = paymentsLogCollection().InsertOne(context.TODO(), entry)
	return
}

//PendingPayment models a pending payment which needs to be fulfilled
type PendingPayment struct {
	ID                  primitive.ObjectID   `bson:"_id" json:"_id" groups:"public"`
	Gateway             PaymentGateway       `bson:"gateway" json:"gateway" groups:"public"`
	Created             time.Time            `bson:"created" json:"created" groups:"public"`
	Anonymous           bool                 `bson:"anonymous"  json:"anonymous" groups:"public"`
	Fullname            string               `bson:"fullname" json:"fullname" groups:"public"`
	Email               string               `bson:"email" json:"email" groups:"public"`
	PhoneCountryCode    string               `bson:"phone_country_code" json:"phone_country_code" groups:"public"`
	PhoneNationalNumber string               `bson:"national_number" json:"national_number" groups:"public"`
	ExtraPhone          string               `bson:"extra_phone" json:"notification_phone" groups:"public"`
	ThankYou            string               `bson:"thank_you" json:"thank_you,omitempty" groups:"public"`
	Comment             string               `bson:"comment" json:"comment" groups:"public"`
	Currency            PaymentCurrency      `bson:"currency" json:"currency" groups:"public"`
	Form                *FilledServiceForm   `bson:"form,omitempty" json:"form" groups:"private"`
	Price               float64              `bson:"price" json:"price" groups:"public"`
	Items               int                  `bson:"items" json:"items" groups:"public"`
	ItemName            string               `bson:"item_name,omitempty" json:"item_name,omitempty" groups:"public"`
	Status              string               `bson:"status" json:"status" groups:"public"`
	Action              PendingPaymentAction `bson:"on_complete" json:"action" groups:"private"`
	RedirectURL         string               `bson:"redirect_url" json:"redirect_url" groups:"protected"`
	PollURL             string               `bson:"poll_url" json:"poll_url" groups:"protected"`
	TargetCreator       primitive.ObjectID   `bson:"creator" json:"creator" groups:"private"`
	UpsertFan           bool                 `bson:"upsert_fan" json:"upsert_fan" groups:"private"` //create fan if it doesnt exist
	TargetFan           primitive.ObjectID   `bson:"fan" json:"fan" groups:"private"`
	TargetCampaign      primitive.ObjectID   `bson:"campaign" json:"campaign_id" groups:"public"`
	Log                 string               `bson:"log" json:"log" groups:"protected"`
	State               string               `bson:"state" json:"state" groups:"protected"`
	UnlockCode          string               `bson:"-" json:"unlock_code,omitempty" groups:"public"`
}

//NewSubscriptionPendingPayment creates a new pending payment for a subscription
func NewSubscriptionPendingPayment(p PendingPayment) (id primitive.ObjectID, err error) {
	p.ID = primitive.NewObjectID()
	p.Action = SubscribeOnPaid

	_, err = pendingPaymentsCollection().InsertOne(context.TODO(), p)
	id = p.ID
	return
}

//NewPayPerViewPendingPayment creates a new pending payment for a pay per view campaign
func NewPayPerViewPendingPayment(p PendingPayment) (id primitive.ObjectID, err error) {
	p.ID = primitive.NewObjectID()
	p.Action = AllowAccessOnPaid

	_, err = pendingPaymentsCollection().InsertOne(context.TODO(), p)
	id = p.ID
	return
}

//NewDonationPendingPayment creates a new pending payment for a "support donation" to a creator
func NewDonationPendingPayment(p PendingPayment) (id primitive.ObjectID, err error) {
	p.ID = primitive.NewObjectID()
	p.Action = AddSupportMessageOnPaid

	_, err = pendingPaymentsCollection().InsertOne(context.TODO(), p)
	id = p.ID
	return
}

//NewVerifyAccountPendingPayment creates a new pending payment for a verification of an account
func NewVerifyAccountPendingPayment(p PendingPayment) (id primitive.ObjectID, err error) {
	p.ID = primitive.NewObjectID()
	p.Action = VerifyAccountOnPaid

	_, err = pendingPaymentsCollection().InsertOne(context.TODO(), p)
	id = p.ID
	return
}

//NewServicePendingPayment creates a new pending payment for a service to a creator
func NewServicePendingPayment(p PendingPayment) (id primitive.ObjectID, err error) {
	p.ID = primitive.NewObjectID()
	p.Action = AllowServiceRequestOnPaid

	_, err = pendingPaymentsCollection().InsertOne(context.TODO(), p)
	id = p.ID
	return
}

//RetrievePendingPayment retrieves a pending payment by id
func RetrievePendingPayment(id primitive.ObjectID) (pending PendingPayment, err error) {
	filter := bson.M{
		"_id": id,
	}

	err = pendingPaymentsCollection().FindOne(context.TODO(), filter).Decode(&pending)
	return
}

//ListRecentPayments lists payments sent or paid for by the user
func (creator User) ListRecentPayments(pageSize int64) (operations []PendingPayment, err error) {
	operations = make([]PendingPayment, 0)
	pageSize = 1000
	var cursor *mongo.Cursor
	filter := bson.M{
		"status": bson.M{"$in": []string{"sent", "paid"}},
		"fan":    creator.ID,
	}
	opts := options.FindOptions{
		Limit: &pageSize,
	}
	opts.SetSort(bson.M{"created": -1})

	cursor, err = pendingPaymentsCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &operations)

	return
}

func GenerateUnlockSignature(paymentID, campaignID, ts string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s!%s!%s!!!nonce#", paymentID, campaignID, ts))))
}

//UpdateStatus updates the status of a pending payment
func (p PendingPayment) UpdateStatus(status, log string) (err error) {
	filter := bson.M{
		"_id": p.ID,
	}

	p.Status = status
	update := bson.M{
		"$set": bson.M{
			"status":   status,
			"poll_url": p.PollURL,
			"log":      p.Log + "\n\n" + log,
		},
	}

	_, err = pendingPaymentsCollection().UpdateOne(context.TODO(), filter, update)
	fmt.Printf("update err: %v", err)
	return
}
