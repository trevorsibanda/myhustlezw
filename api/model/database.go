package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout                   = 10
	connectionStringTemplate         = "mongodb://%s:%s@%s"
	creatorsCollectionName           = "creators"
	otpCollectionName                = "otp"
	UserCredentialsCollectionName    = "user_security"
	walletCollectionName             = "wallet_operations"
	notificationCollectionName       = "notifications"
	campaignsCollectionName          = "campaigns"
	filledCampaignServiceFormsName   = "filled_service_forms"
	supportersCollectionName         = "supporters"
	creatorFilesCollectionName       = "creator_files"
	paymentLogsCollectionName        = "payments_log"
	pendingPaymentsCollectionName    = "pending_payments"
	withdrawalRequestsCollectionName = "withdrawal_requests"
	badUsernamesCollectionName       = "bad_usernames"
)

var (
	mongoConnection *mongo.Client
	mongoDBName     string
	mongoContext    context.Context
	cancelFunc      context.CancelFunc
)

// GetConnection - Retrieves a client to the DocumentDB
func getConnection(username, password, clusterEndpoint string) (*mongo.Client, context.Context, context.CancelFunc) {

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
		panic(err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
		panic(err)
	}

	log.Printf("Connected to MongoDB!")
	return client, ctx, cancel
}

func getDB() *mongo.Client {
	return mongoConnection
}

//WorkingDatabase return an instance of the current active database
func WorkingDatabase() *mongo.Database {
	return getDB().Database(mongoDBName)
}

//creatorsCollection returns the creator users collection
func creatorsCollection() *mongo.Collection {
	return WorkingDatabase().Collection(creatorsCollectionName)
}

func creatorFilesCollection() *mongo.Collection {
	return WorkingDatabase().Collection(creatorFilesCollectionName)
}

//campaignsCollection returns the one time pins collection
func campaignsCollection() *mongo.Collection {
	return WorkingDatabase().Collection(campaignsCollectionName)
}

//otpCollection returns the one time pins collection
func otpCollection() *mongo.Collection {
	return WorkingDatabase().Collection(otpCollectionName)
}

//securityCollection returns the user security details collection
func securityCollection() *mongo.Collection {
	return WorkingDatabase().Collection(UserCredentialsCollectionName)
}

//walletCollection returns the creator wallet collection
func walletCollection() *mongo.Collection {
	return WorkingDatabase().Collection(walletCollectionName)
}

//paymentsLogCollection returns the payment log collection
func paymentsLogCollection() *mongo.Collection {
	return WorkingDatabase().Collection(paymentLogsCollectionName)
}

//badUsernamesCollection returns the disallowed usernames collection
func badUsernamesCollection() *mongo.Collection {
	return WorkingDatabase().Collection(paymentLogsCollectionName)
}

func pendingPaymentsCollection() *mongo.Collection {
	return WorkingDatabase().Collection(pendingPaymentsCollectionName)
}

func withdrawalRequestsCollection() *mongo.Collection {
	return WorkingDatabase().Collection(withdrawalRequestsCollectionName)
}

//supportersCollection returns the creator supporter's collection
func supportersCollection() *mongo.Collection {
	return WorkingDatabase().Collection(supportersCollectionName)
}

//notificationsCollection returns the creator supporter's collection
func notificationsCollection() *mongo.Collection {
	return WorkingDatabase().Collection(notificationCollectionName)
}

// InitDatabaseConnection initializes a connection to a database instance
func InitDatabaseConnection() {
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")
	mongoDBName = os.Getenv("MONGODB_DATABASE")
	mongoConnection, mongoContext, cancelFunc = getConnection(username, password, clusterEndpoint)

}
