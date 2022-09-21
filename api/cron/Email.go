package cron

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/matcornic/hermes/v2"
	"github.com/streadway/amqp"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Configure hermes by setting a theme and your product info
var h hermes.Hermes = hermes.Hermes{
	// Optional Theme
	// Theme: new(Default)
	Product: hermes.Product{
		// Appears in header & footer of e-mails
		Name: "MyHustleZW",
		Link: "https://myhustle.co.zw/",
		// Optional product logo
		Logo:      "https://myhustle.co.zw/assets/img/logo.png",
		Copyright: "Copyright © 2022 Hustle Routine. All rights reserved.",
	},
}

var (
	_emailsPendingQueue                          = "emails:dispatch"
	emailsPendingQueue             *amqp.Queue   = nil
	emailsPendingChannel           *amqp.Channel = nil
	_emailsPriorityQueue                         = "emails:priority"
	emailsPriorityQueue            *amqp.Queue   = nil
	emailsPriorityChannel          *amqp.Channel = nil
	_emailsSmartNotificationsQueue               = "emails:smart_notifications"
)

//Dispatchemail dispatches an email
func ForceDispatchEmailLocal(user model.User, email model.EmailNotification, node string) {
	_, err := sendEmailRequest(user, email)
	if err != nil {
		log.Printf("emails Failed to send email %v", err)
	}
}

type pendingEmailDispatch struct {
	Email model.EmailNotification
	User  model.User
	Node  string
}

//Dispatchemail dispatches a email request
func DispatchEmail(user model.User, email model.EmailNotification, priority bool) {
	if email.Email == "" {
		email.Email = user.Email
	}

	entry := pendingEmailDispatch{
		Email: email,
		User:  user,
		Node:  runnerNode,
	}

	channel := emailsPendingChannel
	queue := emailsPendingQueue
	if priority {
		channel = emailsPriorityChannel
		queue = emailsPriorityQueue
	}

	if bsonEntry, err := bson.Marshal(entry); err != nil {
		log.Printf("emails Failed to bson marshal entry %v", err)
	} else {
		if err := channel.Publish("", queue.Name, false, false, amqp.Publishing{
			ContentType:  "application/bson",
			Body:         bsonEntry,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			//send support email
			SendEngineerSupportEmail(fmt.Sprintf("Failed to publish pending email to queue %v . Force initiating! ", err), email)
			ForceDispatchEmailLocal(user, email, runnerNode)
		}
	}
	return
}

//SendEmailVerification sends an email verification code to the user
func SendEmailVerification(user *model.User, action model.OTPAction) (emailOtp *model.EmailOTP, err error) {
	code := model.GenerateOTPCode(6)
	otpOp := model.OTPOperation{
		ID:         primitive.NewObjectID(),
		Owner:      user.ID,
		Action:     action,
		Code:       code,
		Verified:   false,
		VerifiedAt: time.Now(),
		Message:    fmt.Sprintf("Your %s email verification code is %s", prettyAction(action), code),
		Log:        "",
		ExpiresAt:  time.Now().Add(emailOTPExpiresAfter),
		SentAt:     time.Now(),
	}
	emailOtp = &model.EmailOTP{
		OTPOperation: otpOp,
		Email:        user.Email,
		Provider:     emailProvider,
	}

	notification := model.EmailNotification{
		Template: "email-verification",
		Email:    user.Email,
		Title:    "Action Required: " + prettyAction(action),
	}
	notification.ActionName = "Go to your account"
	notification.ActionURL = "https://myhustle.co.zw/"
	notification.Owner = user.ID
	notification.Type = model.EmailNotificationType
	notification.Message = otpOp.Message

	DispatchEmail(*user, notification, true)

	//save to DB
	err = model.SaveEmailOTP(emailOtp)
	if err != nil {
		log.Printf("Failed to create new email otp entry %v, Reason: %v", emailOtp, err)
	}
	return
}

func sendEmailRequest(user model.User, email model.EmailNotification) (res []byte, err error) {
	from := emailGeneralSenderID
	if email.IsApplicationLog {
		from = emailAppLogSenderID
	}
	subject := strings.TrimSpace(email.Title)
	to := email.Email
	if to == "" {
		to = user.Email
	}

	email.Created = time.Now()
	email.SentAt = time.Now()

	hermesEmail := generateEmailFromTemplate(user, email)
	html, _ := h.GenerateHTML(hermesEmail)
	if clear, err := h.GeneratePlainText(hermesEmail); err == nil {
		email.Message = clear
	}
	message := emailClient.NewMessage(fmt.Sprintf("MyHustleZW <%s>", from), subject, email.Message, to)
	message.SetTracking(true)
	message.SetTrackingOpens(true)
	message.SetTrackingClicks(true)
	message.SetHtml(html)

	response, id, err := emailClient.Send(context.TODO(), message)
	email.Delivered = err == nil
	email.Log = fmt.Sprintf("ID: %v\nResponse: %v\nErr: %v\n\n", id, response, err)

	//log.Printf("send email %v %v %v", response, err, email)
	err = model.AddProcessedNotification(email)
	return
}

//emailsDispatcher looks for queued emails and creates them
func emailsDispatcherRunner() {
	var err error
	log.Printf("emails dispatcher started at %s", time.Now())

	assert(emailsPendingQueue != nil, "emails queue is nil")

	defer emailsPendingChannel.Close()

	msgs, err := emailsPendingChannel.Consume(
		_emailsPendingQueue, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		panic(fmt.Sprintf("emails Failed to consume from queue %v", err))
	}

	for d := range msgs {
		var entry pendingEmailDispatch
		if err := bson.Unmarshal([]byte(d.Body), &entry); err != nil {
			log.Printf("emails Failed to unmarshal message %v", err)
			break
		}
		log.Printf("emails Processing email %s to %s", entry.Email.Template, entry.Email.Email)
		ForceDispatchEmailLocal(entry.User, entry.Email, entry.Node)
		d.Ack(false)
	}

	log.Printf("emails dispatcher stopped at %s", time.Now())

}

//emailsSettlerRunner rnuns at predefined intervals looking for emails which can be settled
func priorityEmailsRunner() {
	var err error
	log.Printf("priority-emails dispatcher started at %s", time.Now())

	assert(emailsPendingQueue != nil, "emails queue is nil")

	defer emailsPendingChannel.Close()

	msgs, err := emailsPriorityChannel.Consume(
		_emailsPriorityQueue, // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		panic(fmt.Sprintf("priority-emails Failed to consume from queue %v", err))
	}

	for d := range msgs {
		var entry pendingEmailDispatch
		if err := bson.Unmarshal([]byte(d.Body), &entry); err != nil {
			log.Printf("priority-emails Failed to unmarshal message %v", err)
			break
		}
		log.Printf("priority-emails Processing email %s to %s", entry.Email.Template, entry.Email.Email)
		ForceDispatchEmailLocal(entry.User, entry.Email, entry.Node)
		d.Ack(false)

	}
	log.Printf("priority-emails dispatcher stopped at %s", time.Now())
}

//emailsDailyNotifications runs once daily and notifies creators who made earnings on a particular day of
// their earning
func smartEmailNotificationsTask(email model.EmailNotification) bool {
	return false
}

//emailsDailyAdministrativeReport sends an administrative report outline of notifications stats in a day
func emailsDailyAdministrativeReport() {

}

func generateEmailFromTemplate(user model.User, email model.EmailNotification) (h hermes.Email) {
	h = hermes.Email{
		Body: hermes.Body{
			Name:      user.Fullname,
			Signature: "Yours Sincerely, ",
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
				"You are receiving this email because you are a member of MyHustle. If you no longer wish to receive these emails, please login to your account and update your notification settings.",
				fmt.Sprintf("This email was sent at %s and intended for %s ", time.Now().Format(time.UnixDate), user.Fullname),
			},
		},
	}

	switch email.Template {
	case "account_login":
		h.Body.Intros = []string{"Account login notification."}
		h.Body.Actions = []hermes.Action{
			{
				Instructions: fmt.Sprintf("You just logged in to your account at %s. If this wasn't you. Please contact support immediately", time.Now().Format("2006-01-02 15:04:05")),
				Button: hermes.Button{
					Color: "#22BC66", // Optional action button color
					Text:  "Contact Support",
					Link:  "https://myhustle.co.zw/?src=email&action=login_alert",
				},
			},
		}
	case "account_signup":
		h.Body.Intros = []string{"MyHustle account created."}
		h.Body.Actions = []hermes.Action{
			{
				Instructions: "You have successfully created your account. Please login with your new details to get started.",
				Button: hermes.Button{
					Color: "#22BC66", // Optional action button color
					Text:  "Login to your account",
					Link:  "https://myhustle.co.zw/creator/dashboard",
				},
			},
		}
	case "notification":
		h.Body.Intros = []string{"Notification from MyHustle."}
		h.Body.Actions = []hermes.Action{
			{
				Instructions: email.Message,
				Button: hermes.Button{
					Color: "#22BC66", // Optional action button color
					Text:  email.ActionName,
					Link:  email.ActionURL,
				},
			},
		}
	case "support_subscriber":
	case "support_pay_per_view":
	case "support_donation":
	case "payment_received":
	case "payment_made":
	case "account_verified":
	case "account_updated":
	case "account_welcome":
	case "new_subscriber_content":
	case "file_processing_complete":
	case "withdrawal_request_received":
	case "withdrawal_request_processed":
	case "application_log":
	case "general":
	default:
		h = hermes.Email{
			Body: hermes.Body{
				Name:   user.Fullname,
				Intros: []string{},
				Actions: []hermes.Action{
					{
						Instructions: email.Message,
						Button: hermes.Button{
							Color: "#22BC66", // Optional action button color
							Text:  email.ActionName,
							Link:  email.ActionURL,
						},
					},
				},
				Outros: []string{
					"Need help, or have questions? Just reply to this email, we'd love to help.",
					"You are receiving this email because you are a member of MyHustle. If you no longer wish to receive these emails, please login to your account and update your notification settings.",
				},
			},
		}

	}

	if len(email.Dictionary) > 0 {
		h.Body.Dictionary = toHermesDictionary(email.Dictionary)
	}
	return
}

func toHermesDictionary(sd map[string]string) (res []hermes.Entry) {
	for k, v := range sd {
		res = append(res, hermes.Entry{Key: k, Value: v})
	}
	return
}
