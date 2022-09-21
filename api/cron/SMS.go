package cron

import (
	"fmt"
	"log"
	"time"

	"github.com/trevorsibanda/myhustlezw/api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	smsOTPExpiresAfter   = time.Hour * 1
	emailOTPExpiresAfter = time.Hour * 6
	smsGateway           = "twilio"
	emailProvider        = "sendgrid"
	channelsProvider     = "pusher"
)

func sendSMSRequest(notification model.SMSNotification) {
	notification.Type = model.SMSNotificationType
	notification.Created = time.Now()
	notification.SentAt = time.Now()
	response, exception, err := smsClient.SendSMS(smsSenderID, notification.PhoneNumber, notification.Message, "", "")
	if exception != nil || err != nil {
		notification.Delivered = false
		notification.Log = fmt.Sprintf("Failed to send SMSOTP to %v with errors:\n Exception: %v\nError: %v\nResponse: %v", notification, exception, err, response)
	} else {
		notification.Delivered = true
		notification.Log = fmt.Sprintf("SMS Delivery:\nResponse: %v\nException: %v\nError: %v", response, exception, err)
	}

	log.Printf("send sms %v %v %v", response, err, notification)
	err = model.AddProcessedNotification(notification)
	return
}

func prettyAction(action model.OTPAction) string {
	switch action {
	case model.AccountSignup:
		return "account signup"
	case model.ResetPassword:
		return "password reset"
	case model.ChangeEmail:
		return "email change"
	case model.ChangePhone:
		return "phone number change"
	default:
		return "Unknown"
	}
}

//SendSMSVerification sends an SMS verification code to the specified number
func SendSMSVerification(user *model.User, action model.OTPAction) (smsOtp *model.SMSOTP, err error) {
	code := model.GenerateOTPCode(6)
	otpOp := model.OTPOperation{
		ID:        primitive.NewObjectID(),
		Owner:     user.ID,
		Action:    action,
		Code:      code,
		Verified:  false,
		Message:   fmt.Sprintf("Your %s verification code is %s.\nNever share this code with anyone.", prettyAction(action), code),
		Log:       "",
		ExpiresAt: time.Now().Add(smsOTPExpiresAfter),
		SentAt:    time.Now(),
	}

	smsOtp = &model.SMSOTP{
		OTPOperation: otpOp,
		PhoneNumber:  user.PhoneNumber,
		Gateway:      smsGateway,
	}

	notification := &model.SMSNotification{
		PhoneNumber: smsOtp.PhoneNumber,
	}

	notification.Message = smsOtp.Message
	notification.Type = model.SMSNotificationType
	notification.Created = time.Now()
	notification.SentAt = time.Now()
	notification.Priority = true

	err = SendSMSNotification(notification)
	if err != nil {
		smsOtp.Delivered = false
		smsOtp.Log = fmt.Sprintf("Failed to send SMSOTP to %v with errors: %v", smsOtp, err)
	} else {
		smsOtp.Delivered = true
	}

	//save to DB
	err = model.SaveSMSOTP(smsOtp)
	return
}

//SendSMSNotification is used to send a general SMS nootification message
func SendSMSNotification(notification *model.SMSNotification) (err error) {
	notification.Type = model.SMSNotificationType
	notification.Created = time.Now()
	notification.SentAt = time.Now()
	response, exception, err := smsClient.SendSMS(smsSenderID, notification.PhoneNumber, notification.Message, "", "")
	if exception != nil || err != nil {
		notification.Delivered = false
		notification.Log = fmt.Sprintf("Failed to send SMSOTP to %v with errors:\n Exception: %v\nError: %v\nResponse: %v", notification, exception, err, response)
	} else {

		notification.Delivered = true
		notification.Log = fmt.Sprintf("SMS Delivery:\nResponse: %v\nException: %v\nError: %v", response, exception, err)
	}

	err = model.AddProcessedNotification(notification)
	return
}
