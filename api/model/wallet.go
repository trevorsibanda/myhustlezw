package model

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PayoutMethod string

const (
	NoPayoutMethod    PayoutMethod = "none"
	PayoutBank        PayoutMethod = "bank"
	PayoutMobileMoney PayoutMethod = "mobile_money"
)

//WalletOperation models an operation on a user's wallet
//This includes credit and debit operations
type WalletOperation struct {
	ID          primitive.ObjectID `bson:"_id"  json:"_id" groups:"public"`
	Creator     primitive.ObjectID `bson:"creator_id"  json:"creator_id" groups:"public"`
	Supporter   primitive.ObjectID `bson:"supporter_id" json:"supporter_id" groups:"public"`
	Operation   string             `bson:"operation"  json:"operation" groups:"public"`
	Currency    string             `bson:"currency" json:"currency" groups:"public"`
	Gateway     string             `bson:"gateway"  json:"gateway" groups:"public"`
	Created     time.Time          `bson:"created"  json:"created" groups:"public"`
	ReleaseAt   time.Time          `bson:"release" json:"release" groups:"public"`
	LastUpdated time.Time          `bson:"last_updated" json:"last_updated" groups:"public"`
	Amount      float64            `bson:"amount"  json:"amount" groups:"public"`
	Comment     string             `bson:"comment"  json:"comment" groups:"protected"`
}

//RetrieveEscrowReadyWalletOps retrieves wallet ops in escrow that are ready to be depositted to a user's account
func RetrieveEscrowReadyWalletOps(limit int64) (ops []WalletOperation, err error) {
	ops = make([]WalletOperation, 0)
	filter := bson.M{
		"operation": "credit_escrow",
		"release": bson.M{
			"$lte": time.Now(),
		},
	}

	opts := options.FindOptions{
		Limit: &limit,
	}

	opts.SetSort(bson.M{"created": -1})

	cursor, err := walletCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	err = cursor.All(context.TODO(), &ops)
	return
}

//RequestWithdrawal marks all ready funds as pending withdrawal
//Run a WalletSummary after this
func RequestWithdrawal(creator User, currency PaymentCurrency, comment string) (err error) {
	filter := bson.M{
		"operation":  "credit_available",
		"creator_id": creator.ID,
		"currency":   currency,
	}

	_, err = walletCollection().UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{
			"comment":   comment,
			"operation": "debit_pending_withdrawal",
		},
	})
	return
}

//Dispute disputes a payment with the provided comment
func (op *WalletOperation) Dispute(comment string) (err error) {
	op.Operation = "credit_disputed"
	op.LastUpdated = time.Now()
	op.Comment = fmt.Sprintf("%v\nMoved wallet op from %v to disputed\nProvided comment: %v\n\n", op.Comment, op.Operation, comment)
	err = walletCollection().FindOneAndUpdate(context.TODO(), bson.M{
		"_id": op.ID,
	}, bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"operation":    op.Operation,
			"comment":      op.Comment,
		},
	}).Decode(op)
	return
}

//ApproveRefund processes a disputed transaction and refunds it
func (op *WalletOperation) ApproveRefund(comment string) (err error) {
	if op.Operation != "credit_disputed" {
		err = fmt.Errorf("Can only refund disputed transactions. Status is %v", op.Operation)
		return
	}
	op.Operation = "debit_refunded"
	op.LastUpdated = time.Now()
	op.Comment = fmt.Sprintf("%v\nMoved wallet op from %v to refunded\nProvided comment: %v\n\n", op.Comment, op.Operation, comment)
	err = walletCollection().FindOneAndUpdate(context.TODO(), bson.M{
		"_id": op.ID,
	}, bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"operation":    op.Operation,
			"comment":      op.Comment,
		},
	}).Decode(op)
	return
}

//Available marks a transaction as ready to withdraw
func (op *WalletOperation) Available(comment string) (err error) {
	if op.Operation != "credit_escrow" {
		err = fmt.Errorf("Can only move escrow funds to available. Status is %v", op.Operation)
		return
	}
	op.Operation = "credit_available"
	op.LastUpdated = time.Now()
	op.Comment = fmt.Sprintf("%v\nMoved wallet op from %v to updated\nProvided comment: %v\n\n", op.Comment, op.Operation, comment)
	err = walletCollection().FindOneAndUpdate(context.TODO(), bson.M{
		"_id": op.ID,
	}, bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"operation":    op.Operation,
			"comment":      op.Comment,
		},
	}).Decode(op)
	return
}

//Withdrawn marks a payment as withdrawn
func (op *WalletOperation) Withdrawn(comment string) (err error) {
	if op.Operation != "debit_pending_withdrawal" {
		err = fmt.Errorf("Can only withdraw funds that are marked as pending withdrawal. Status is %v", op.Operation)
		return
	}
	op.Operation = "debit_withdrawn"
	op.LastUpdated = time.Now()
	op.Comment = fmt.Sprintf("%v\nMoved wallet op from %v to withdrawn\nProvided comment: %v\n\n", op.Comment, op.Operation, comment)
	err = walletCollection().FindOneAndUpdate(context.TODO(), bson.M{
		"_id": op.ID,
	}, bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"operation":    op.Operation,
			"comment":      op.Comment,
		},
	}).Decode(op)
	return
}

//CreatorWalletSummary provides a summary of a creator's wallet operations
type CreatorWalletSummary struct {
	Currency          string  `groups:"private" json:"currency"`
	Escrow            float64 `groups:"private" json:"escrow"`
	Available         float64 `groups:"private" json:"available"`
	Withdrawn         float64 `groups:"private" json:"withdrawn"`
	PendingWithdrawal float64 `groups:"private" json:"pending_withdrawal"`
	Disputed          float64 `groups:"private" json:"disputed"`
	Refunded          float64 `groups:"private" json:"refunded"`
	InAccurate        bool    `groups:"private" json:"inaccurate"`
}

func (summary CreatorWalletSummary) ApproveWithdrawalURL(creator primitive.ObjectID) string {
	return fmt.Sprintf("%s://%s/api/v1/public/_admin/_/wallet/withdrawal/approve/%s/%f/%s", myhustleScheme, myhustleEndpoint, summary.Currency, summary.PendingWithdrawal, creator.Hex())
}

func (summary *CreatorWalletSummary) ApprovePendingWithdrawal(creatorID primitive.ObjectID) (count int64, err error) {

	filter := bson.M{
		"operation":  "debit_pending_withdrawal",
		"creator_id": creatorID,
		"currency":   summary.Currency,
	}

	result, err1 := walletCollection().UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
			"operation":    "debit_withdrawn",
		},
	})
	err = err1
	if result != nil {
		count = result.ModifiedCount
	}
	return
}

//CreatorPayoutFields are fields for a payout
type CreatorPayoutFields struct {
	Method            PayoutMethod `bson:"method" json:"method" groups:"private"`
	BankName          string       `bson:"bank_name" json:"bankname" groups:"private"`
	BankBranch        string       `bson:"bank_branch" json:"bankbranch" groups:"private"`
	BankAccountName   string       `bson:"bank_account_name" json:"bankaccountname" groups:"private"`
	BankAccountNumber string       `bson:"bank_account_number" json:"bankaccountnumber" groups:"private"`
	PhoneNumber       string       `bson:"phone_number" json:"phonenumber" groups:"private"`
}

//CreatorPayoutDetails models details we use for creator payouts
type CreatorPayoutDetails struct {
	USD CreatorPayoutFields `bson:"usd"`
	ZWL CreatorPayoutFields `bson:"zwl"`
}

func (payout CreatorPayoutDetails) For(currency string) (target string) {
	fields := payout.USD
	if strings.ToLower(currency) == "zwl" {
		fields = payout.ZWL
	}

	if fields.Method == PayoutBank {
		target = fmt.Sprintf("Bank:%s\nBranch: %s\nAccount Name: %s\nAccount Number: %s", fields.BankName, fields.BankBranch, fields.BankAccountName, fields.BankAccountNumber)
	} else if fields.Method == PayoutMobileMoney {
		target = fmt.Sprintf("Mobile money\n%s", fields.PhoneNumber)
	} else {
		target = fmt.Sprintf("Unknown payout - %s", fields.Method)
	}
	return
}

//UpdatePayoutDetails updates the payout details
func (creator User) UpdatePayoutDetails(currency PaymentCurrency, fields CreatorPayoutFields) (err error) {
	details := creator.PayoutDetails
	switch currency {
	case ZWL:
		details.ZWL = fields
	case USD:
		details.USD = fields
	}

	filter := bson.M{
		"_id": creator.ID,
	}

	doc := bson.M{
		"$set": bson.M{
			"payout": details,
		},
	}

	err = creatorsCollection().FindOneAndUpdate(context.TODO(), filter, doc).Err()
	return
}

//WalletSummary iterates over all wallet operations and generates a summary
func (creator User) WalletSummary(currency string) (summary CreatorWalletSummary, err error) {
	summary = CreatorWalletSummary{
		Currency: currency,
	}
	filter := bson.M{
		"currency":   currency,
		"creator_id": creator.ID,
	}
	cursor, err := walletCollection().Find(context.TODO(), filter)
	if err != nil {
		return
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var op WalletOperation
		err := cursor.Decode(&op)
		if err != nil {
			summary.InAccurate = true
			log.Printf("Encountered inaccurate wallet operation %v %v %v", creator, currency, err)
			continue
		}
		switch op.Operation {
		case "credit_escrow":
			summary.Escrow += op.Amount
		case "credit_available":
			summary.Available += op.Amount
		case "debit_pending_withdrawal":
			summary.PendingWithdrawal += op.Amount
		case "debit_withdrawn":
			summary.Withdrawn += op.Amount
		case "credit_disputed":
			summary.Disputed += op.Amount
		case "debit_refunded":
			summary.Refunded += op.Amount
		}
	}

	return
}

//CrebitEscrow creates a new escrow deposit operation
func (creator User) CreditEscrow(amount float64, gateway, currency, comment string, fanID primitive.ObjectID) (op WalletOperation, err error) {

	op = WalletOperation{
		ID:          primitive.NewObjectID(),
		Operation:   "credit_escrow",
		Creator:     creator.ID,
		Gateway:     gateway,
		Currency:    currency,
		Amount:      amount,
		Comment:     comment,
		Created:     time.Now(),
		ReleaseAt:   time.Now().Add(72 * time.Hour), //hold in escrow for 72hours
		LastUpdated: time.Now(),
	}

	if fanID.IsZero() {
		op.Supporter = fanID
	}

	_, err = walletCollection().InsertOne(context.TODO(), op)
	return
}

//ListRecentWalletOperations returns a list of recent wallet activity
func (creator User) ListRecentWalletOperations(pageSize int64) (operations []WalletOperation, err error) {
	operations = make([]WalletOperation, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"creator_id": creator.ID,
	}
	opts := options.FindOptions{
		Limit: &pageSize,
	}
	opts.SetSort(bson.M{"created": -1})

	cursor, err = walletCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &operations)

	return
}

//GetWalletOp returns a wallet entry
func (creator User) GetWalletOp(id primitive.ObjectID) (op WalletOperation, err error) {
	filter := bson.M{
		"creator_id": creator.ID,
		"_id":        id,
	}

	err = walletCollection().FindOne(context.TODO(), filter).Decode(&op)
	return
}
