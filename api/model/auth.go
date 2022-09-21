package model

import (
	"context"
	"fmt"
	"math/rand"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	smsOTPExpiresAfter   = time.Hour * 1
	emailOTPExpiresAfter = time.Hour * 6
	smsGateway           = "twilio"
	emailProvider        = "sendgrid"
	channelsProvider     = "pusher"
)

//OTPAction is the action driving an OTP request
type OTPAction string

const (
	//AccountSignup is an otp action which when completed, verifies the account during signup
	AccountSignup OTPAction = "account_signup"
	//ResetPassword is an otp action which when completed allows the user to specify a new password
	ResetPassword OTPAction = "reset_password"
	//ChangeEmail is an otp action which when completed allows the user to change their email address
	ChangeEmail OTPAction = "change_email"
	//ChangePhone is an otp action which when completed allows the user to change their phone number
	ChangePhone OTPAction = "change_phone"
)

//OTPOperation is a one time pin operation
type OTPOperation struct {
	ID         primitive.ObjectID `bson:"_id"  json:"_id" groups:"public"`
	Owner      primitive.ObjectID `bson:"owner_id"  json:"owner_id" groups:"private"`
	Code       string             `bson:"code"  json:"code" groups:"private"`
	Action     OTPAction          `bson:"action"  json:"action" groups:"private"`
	Verified   bool               `bson:"verified"  json:"verified" groups:"public"`
	Delivered  bool               `bson:"delivered"  json:"delivered" groups:"public"`
	SentAt     time.Time          `bson:"sent_at"  json:"sent_at" groups:"public"`
	ExpiresAt  time.Time          `bson:"expires_at"  json:"expires_at" groups:"public"`
	VerifiedAt time.Time          `bson:"verified_at"  json:"verified_at" groups:"public"`
	Message    string             `bson:"message"  json:"message" groups:"protected"`
	Log        string             `bson:"log"  json:"log" groups:"protected"`
}

//SMSOTP represents an auth operation where an sms is sent to a user's
//mobile device.
type SMSOTP struct {
	OTPOperation `bson:",inline"`
	PhoneNumber  string `bson:"phone_number"  json:"phone_number" groups:"private"`
	Gateway      string `bson:"gateway"  json:"gateway" groups:"private"`
}

//EmailOTP models an auth operation where a verification code is sent to a user's
//email address.
type EmailOTP struct {
	OTPOperation `bson:",inline"`
	Email        string `bson:"email"  json:"email" groups:"private"`
	Provider     string `bson:"provider"  json:"provider" groups:"private"`
}

//UserCredentials models security info we store about the user
type UserCredentials struct {
	Owner                  primitive.ObjectID `bson:"owner_id"  json:"owner_id" groups:"private"`
	Password               string             `bson:"password"  json:"password" groups:"protected"`
	MobileConfirmationCode string             `bson:"mobile_confirmation"  json:"mobile_confirmation" groups:"protected"`
	EmailConfirmationCode  string             `bson:"email_confirmation"  json:"email_confirmation" groups:"protected"`
}

//RetrieveUserCredentials retrieves the user's security credentials
func RetrieveUserCredentials(userID primitive.ObjectID) (cred UserCredentials, err error) {
	filter := bson.M{
		"owner_id": userID,
	}
	err = securityCollection().FindOne(context.TODO(), filter).Decode(&cred)
	return
}

//CanSendVerification checks if the user can send another OTP message
//for the given action. 3 actions are allowed per hour.
func CanSendVerification(user *User, action OTPAction) (b bool, err error) {
	var count int64 = 9
	filter := bson.M{
		"owner_id": user.ID,
		"action":   action,
		"verified": false,
		"sent_at":  bson.M{"$gt": time.Now().Add(-time.Minute * 60)},
	}
	count, err = otpCollection().CountDocuments(context.TODO(), filter)

	b = (count < 3)
	return
}

//VerifySMSCode verifies a code received against a OTP
func VerifySMSCode(userID primitive.ObjectID, phone, code string, action OTPAction) (otp SMSOTP, err error) {
	filter := bson.M{
		"owner_id":     userID,
		"phone_number": phone,
		"code":         code,
		"action":       action,
		"verified":     false,
	}
	err = otpCollection().FindOne(context.TODO(), filter).Decode(&otp)
	if err == nil {
		if !otp.ExpiresAt.After(time.Now()) {
			err = fmt.Errorf("The verification code has expired")
		} else {
			otp.Verified = true
			otp.VerifiedAt = time.Now()
			err = otpCollection().FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": bson.M{"verified": true, "verified_ever": true, "verified_at": time.Now()}}).Err()
		}
	}
	return
}

//VerifyEmailCode verifies a code received against a OTP
func VerifyEmailCode(userID primitive.ObjectID, email, code string, action OTPAction) (otp SMSOTP, err error) {
	filter := bson.M{
		"owner_id": userID,
		"email":    email,
		"code":     code,
		"action":   action,
		"verified": false,
	}
	err = otpCollection().FindOne(context.TODO(), filter).Decode(&otp)
	if err == nil {
		if !otp.ExpiresAt.After(time.Now()) {
			err = fmt.Errorf("The verification code has expired")
		} else {
			otp.Verified = true
			otp.VerifiedAt = time.Now()
			err = otpCollection().FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": bson.M{"verified": true, "verified_ever": true, "verified_at": time.Now()}}).Err()
		}
	}
	return
}

//ValidatePasswords compares two user passwords
func ValidatePasswords(password1, password2 string) bool {
	return compareUserPasswords(password1, password2)
}

//newUserCredentials adds security details for a new user.
func newUserCredentials(user *User, password string) (sec *UserCredentials, err error) {
	sec = &UserCredentials{
		Owner:                  user.ID,
		Password:               password,
		MobileConfirmationCode: GenerateOTPCode(5),
		EmailConfirmationCode:  GenerateOTPCode(6),
	}
	_, err = securityCollection().InsertOne(context.TODO(), sec)
	return
}

func GenerateOTPCode(n int) string {
	letterBytes := "0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SaveSMSOTP(otp *SMSOTP) (err error) {
	_, err = otpCollection().InsertOne(context.TODO(), otp)
	return
}

func SaveEmailOTP(otp *EmailOTP) (err error) {
	_, err = otpCollection().InsertOne(context.TODO(), otp)
	return
}

//compareUserPasswords compares passwords
//TODO: add encryption
func compareUserPasswords(password1 string, password2 string) bool {
	return password1 == password2
}
