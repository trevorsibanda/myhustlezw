package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/pusher/pusher-http-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	smsClient           *gotwilio.Twilio
	emailClient         *mailgun.MailgunImpl
	channelsClient      *pusher.Client
	smsAccountSID       string
	smsAuthToken        string
	smsSenderID         string
	emailSendgridToken  string
	pusherConnectionURL string

	emailGeneralSenderID   string
	emailAppLogSenderID    string
	financeDepartmentEmail string
	financeDepartmentPhone string
)

//InitializeOTPClients creates a SMS client and email client for use in sending OTP codes
func InitializeNotificationClients() {
	smsAccountSID = os.Getenv("SMS_TWILIO_ACCOUNT_SID") // "ACedd0e47745dbfefc6cc7c40dd48c5952"
	smsAuthToken = os.Getenv("SMS_TWILIO_AUTH_TOKEN")   //  "30950aa65c8ce1d54ba41252b4365f2f"
	smsSenderID = os.Getenv("SMS_TWILIO_PHONE_NUMBER")  //"+12812385674"
	emailSendgridToken = os.Getenv("SENDGRID_API_KEY")
	smsClient = gotwilio.NewTwilioClient(smsAccountSID, smsAuthToken)
	emailClient = mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_PRIVATE_KEY"))

	emailGeneralSenderID = os.Getenv("EMAIL_GENERAL_SENDER_ID")
	financeDepartmentEmail = os.Getenv("EMAIL_FINANCE_DEPARTMENT")
	financeDepartmentPhone = os.Getenv("PHONE_FINANCE_DEPARTMENT")
	emailAppLogSenderID = os.Getenv("EMAIL_APPLOG_SENDER_ID")

	pusherConnectionURL = os.Getenv("PUSHER_CONNECTION_URL")
	var err error
	if channelsClient, err = pusher.ClientFromURL(pusherConnectionURL); err != nil {
		panic(fmt.Sprintf("Failed to create pusher client from connection URL - %s", pusherConnectionURL))
	}
}

func SendFinanceDepartmentNotification(creator model.User, summary model.CreatorWalletSummary, data map[string]interface{}) {
	creatorSMS := &model.SMSNotification{
		PhoneNumber: financeDepartmentPhone,
	}

	creatorSMS.Message = fmt.Sprintf("New withdrawal request  for $%s by @%s. Details sent to email.", Format(summary.Currency, summary.PendingWithdrawal), creator.Username)
	creatorSMS.Owner = creator.ID

	SendSMSNotification(creatorSMS)

	creatorEmail := model.EmailNotification{
		Email:      financeDepartmentEmail,
		Title:      fmt.Sprintf("Withdrawal request by %s @%s of $%s %d", creator.Email, creator.Username, summary.Currency, summary.PendingWithdrawal),
		Template:   "transaction",
		ActionURL:  summary.ApproveWithdrawalURL(creator.ID),
		ActionName: "Approve withdrawal and notify of success",
	}
	creatorEmail.Message = fmt.Sprintf("New withdrawal request  for $%s by @%s.", Format(summary.Currency, summary.PendingWithdrawal), creator.Username)
	prettyData, _ := json.MarshalIndent(data, " ", "\t")
	creatorEmail.Dictionary = map[string]string{
		"Currency":              string(summary.Currency),
		"Processing Withdrawal": Format(summary.Currency, summary.PendingWithdrawal),
		"Amount to Payout":      Format(summary.Currency, summary.PendingWithdrawal*0.7),
		"MyHustle Fees":         Format(summary.Currency, summary.PendingWithdrawal*0.3),
		"Available":             Format(summary.Currency, summary.Available),
		"Disputed":              Format(summary.Currency, summary.Disputed),
		"Payout To":             creator.PayoutDetails.For(summary.Currency),
		"Data":                  string(prettyData),
	}

	DispatchEmail(creator, creatorEmail, true)
}

//SendEngineerSupportEmail sends an email to the engineer support team
func SendEngineerSupportEmail(title string, data interface{}) {
	log.Printf("TODO: Log engineer support email with title %s and data %v", title, data)
}

//SendWalletChangeEmail sends an email when a wallet change event occurs
func SendWalletChangeEmail(title string, message string, creator model.User) {

	creatorEmail := model.EmailNotification{
		Email:      creator.Email,
		Title:      "Funds update: " + title,
		Template:   "creator_withdrawal",
		ActionURL:  creator.URL(),
		ActionName: "View account",
	}
	creatorEmail.Message = message

	DispatchEmail(creator, creatorEmail, true)

}

//SendWalletWithdrawalRequest sends notifications of withdrawals
func SendWalletWithdrawalRequest(creator model.User, summary model.CreatorWalletSummary) (err error) {

	creatorSMS := &model.SMSNotification{
		PhoneNumber: creator.PhoneNumber,
	}

	creatorSMS.Message = fmt.Sprintf("Hi @%s,\nWe received your withdrawal request. More details sent to your email.", creator.Username)
	creatorSMS.Owner = creator.ID

	SendSMSNotification(creatorSMS)

	creatorEmail := model.EmailNotification{
		Email:      creator.Email,
		Title:      "Withdrawal request confirmation.",
		Template:   "creator_withdrawal",
		ActionURL:  creator.URL(),
		ActionName: "View account",
	}
	creatorEmail.Message = fmt.Sprintf("Your request to withdraw %s has been received. You will receive a notification once it is processed.", Format(summary.Currency, summary.PendingWithdrawal))
	creatorEmail.Dictionary = map[string]string{
		"Currency":         string(summary.Currency),
		"MyHustle Fees":    Format(summary.Currency, summary.PendingWithdrawal*0.3),
		"You will receive": Format(summary.Currency, summary.PendingWithdrawal*0.7),
		"Payout To":        creator.PayoutDetails.For(summary.Currency),
	}

	if summary.Available > 0 {
		creatorEmail.Dictionary["Available"] = Format(summary.Currency, summary.Available)
	}

	if summary.Disputed > 0 {
		creatorEmail.Dictionary["Disputed"] = Format(summary.Currency, summary.Disputed)
	}

	DispatchEmail(creator, creatorEmail, true)

	//send email to finance dept
	SendFinanceDepartmentNotification(creator, summary, map[string]interface{}{
		"creator": creator,
		"summary": summary,
	})

	return
}

func SendWalletWithdrawalApproved(creator model.User, summary model.CreatorWalletSummary, header string) (err error) {

	creatorSMS := &model.SMSNotification{
		PhoneNumber: creator.PhoneNumber,
	}

	creatorSMS.Message = fmt.Sprintf("Hi @%s,\nYour withdrawal request has been processed. More details sent to your email.", creator.Username)
	creatorSMS.Owner = creator.ID

	SendSMSNotification(creatorSMS)

	creatorEmail := model.EmailNotification{
		Email:      creator.Email,
		Title:      "We just processed your withdrawal.",
		Template:   "creator_withdrawal",
		ActionURL:  creator.URL(),
		ActionName: "View my account",
	}
	payout := Format(summary.Currency, summary.PendingWithdrawal*0.7)
	fees := Format(summary.Currency, summary.PendingWithdrawal*0.3)
	creatorEmail.Message = fmt.Sprintf("%s\nWe have processed your withdrawal of %s. You will receive %s in your bank account and %s was retained by MyHustle as fees. Please expect the funds to be in your nominated payout bank/mobile money account within 3-4 business days.", header, Format(summary.Currency, summary.PendingWithdrawal), payout, fees)

	DispatchEmail(creator, creatorEmail, true)

	data := bson.M{
		"creator": creator,
		"summary": summary,
	}

	//send email to finance dept
	financeEmail := model.EmailNotification{
		Email:    financeDepartmentEmail,
		Title:    fmt.Sprintf("Withdrawal by @%s of $%s marked as processed!", creator.Username, payout),
		Template: "transaction",
	}
	financeEmail.Message = fmt.Sprintf("Marked withdrawal request of %s by %s %s @%s as processed. If there is any error between now and then, please contact the creator on their contact number or email to resolve.", payout, creator.Email, creator.PhoneNumber, creator.Username)
	prettyData, _ := json.MarshalIndent(data, " ", "\t")
	financeEmail.Dictionary = map[string]string{
		"Data": string(prettyData),
	}
	financeEmail.Dictionary = map[string]string{
		"Currency":         string(summary.Currency),
		"Amount withdrawn": Format(summary.Currency, summary.PendingWithdrawal),
	}
	DispatchEmail(creator, financeEmail, true)

	return
}

//SendPayoutDetailsChangedNotification sends payout details
func SendPayoutDetailsChangedNotification(creator model.User, currency model.PaymentCurrency, details model.CreatorPayoutFields) (err error) {
	email := model.EmailNotification{
		Title:      "Your bank payout details have been changed",
		Template:   "account_change",
		ActionURL:  creator.URL(),
		ActionName: "Go to your account",
		Email:      creator.Email,
	}
	email.Owner = creator.ID
	email.Message = fmt.Sprintf("A change has been made to your account.\n\n\nYou changed your %s payout details.", currency)
	DispatchEmail(creator, email, true)
	return
}

//SendDebitAccountFailedEmail sends an email to creator and admin that it failed to debit an account
func SendDebitAccountFailedEmail(creator model.User, payment model.PendingPayment) {
	email := model.EmailNotification{
		Title:      "You received money but we failed to debit your account.",
		Template:   "transaction",
		ActionURL:  creator.URL(),
		ActionName: "Contact support",
		Email:      creator.Email,
	}
	email.Owner = creator.ID
	email.Message = fmt.Sprintf(
		"You received $%s %f through %s with payment ID %s at %s but we failed to debit your account. Please contact support with the details in this email.",
		payment.Currency, payment.Price, payment.Gateway, payment.ID, payment.Created.Format("2006-01-02 15:04:05"))
	DispatchEmail(creator, email, true)
}

//SendChannelMessage sends a message to a given channel
func SendChannelMessage(channel, event string, data interface{}) (err error) {
	err = channelsClient.Trigger(channel, event, data)
	return
}

func SendNewSubscriberNotification(supporter model.User, creator model.User, payment model.PendingPayment, walletOp model.WalletOperation, support model.CreatorSupport) (err error) {
	//notify supporter of success paid
	var amount string
	switch payment.Currency {
	case model.USD:
		amount = FormatAsUSD(payment.Price)
	case model.ZWL:
		amount = FormatAsZWL(payment.Price)
	}
	SMSmsg := fmt.Sprintf("You have paid %s and subscribed to @%s", amount, creator.Username)
	emailMsg := fmt.Sprintf("Hey %s,\n%s", payment.Fullname, SMSmsg)
	notifyPayerSMS := model.SMSNotification{
		PhoneNumber: support.PhoneNumber,
	}
	notifyPayerSMS.Message = SMSmsg
	notifyPayerSMS.Owner = creator.ID

	notifyPayerEmail := model.EmailNotification{
		Title:      fmt.Sprintf("Paid %s and subscribed to @%s", amount, creator.Username),
		Template:   "support_message",
		ActionURL:  creator.URL(),
		ActionName: fmt.Sprintf("View %s's page", creator.Username),
		Dictionary: map[string]string{
			"Amount Paid": Format(string(payment.Currency), payment.Price),
			"Gateway":     string(payment.Gateway),
			"Payment ID":  payment.ID.Hex(),
		},
	}
	notifyPayerEmail.Email = support.Email
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = creator.ID

	_ = SendSMSNotification(&notifyPayerSMS)

	DispatchEmail(supporter, notifyPayerEmail, false)

	if creator.Notifications.NewSubscriber {
		//TODO: smart notifications
		emailMsg = fmt.Sprintf("Yay!,\n%s paid %s and subscribed to your account.\n\n You can view this payment in detail in your dashboard", support.DisplayName, amount)
		notifyPayeeEmail := model.EmailNotification{
			Title:      fmt.Sprintf("%s bought you %s", support.DisplayName, util.PluralOf(payment.Items, payment.ItemName)),
			Template:   "support_message_creator",
			ActionURL:  creator.URL(),
			ActionName: "Go to your dashboard",
			Dictionary: map[string]string{
				"Amount": Format(string(payment.Currency), payment.Price),
			},
		}
		notifyPayeeEmail.Email = creator.Email
		notifyPayeeEmail.Message = emailMsg
		notifyPayeeEmail.Owner = creator.ID
		DispatchEmail(creator, notifyPayeeEmail, false)
	}
	return
}

//SendIdentityVerifiedNotification sends an email, sms to the user that their identity has been verified
func SendIdentityVerifiedNotification(user *model.User) {
	SMSmsg := "Your identity has been verified. You can now access all MyHustle features!"
	creatorSMS := &model.SMSNotification{
		PhoneNumber: user.PhoneNumber,
	}

	creatorSMS.Message = fmt.Sprintf("%s,\nWe received your payment, your account is now verified.", user.Fullname)
	creatorSMS.Owner = user.ID

	SendSMSNotification(creatorSMS)
	emailMsg := SMSmsg

	notifyPayerEmail := model.EmailNotification{
		Title:      "Your identity is now verified",
		Template:   "support_message",
		ActionURL:  user.URL(),
		ActionName: "View your page",
		Dictionary: map[string]string{
			"Fullname":     user.Fullname,
			"Username":     fmt.Sprintf("@%s", user.Username),
			"Email":        user.Email,
			"Phone number": user.PhoneNumber,
			"Your page":    user.URL(),
		},
	}
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = user.ID

	DispatchEmail(*user, notifyPayerEmail, true)
}

func ScheduleNotification(notification *model.Notification, user *model.User) {

}

func PushNotifyFansNewPost(campaign *model.Campaign) {

}

func PushNotifyCreatorPostReady(campaign *model.Campaign) {

}

func PushNotifyCreatorNewSupporter(creator *model.User) {

}

func PushNotifyFansSubscriptionsChange(creator *model.User) {

}

func PushNotifyNewCreatorLogin(creator model.User, session sessions.VisitorSession, sec *model.UserCredentials) {
	mail := model.EmailNotification{
		Email:      creator.Email,
		Title:      "New login to your account",
		Template:   "account_login",
		ActionURL:  creator.URL(),
		ActionName: "View account",
		Dictionary: session.Info.Dictionary(),
	}
	mail.Priority = true
	mail.Owner = creator.ID
	DispatchEmail(creator, mail, true)
}

func PushTestAlert(user *model.User) (err error) {
	// Send Notification
	s := &webpush.Subscription{}
	log.Println(user.WebNotification)
	resp, err := webpush.SendNotification([]byte("Fans Only: Thursday gym workout session"), &user.WebNotification.Mobile, &webpush.Options{
		Subscriber:      user.Email,
		VAPIDPublicKey:  os.Getenv("WEBPUSH_VAPID_PUBLIC_KEY"),
		VAPIDPrivateKey: os.Getenv("WEBPUSH_VAPID_PRIVATE_KEY"),
		TTL:             30,
	})
	if err != nil {
		log.Println(resp, err, user.WebNotification, s, os.Getenv("WEBPUSH_VAPID_PUBLIC_KEY"), "private:", os.Getenv("WEBPUSH_VAPID_PRIVATE_KEY"))
	}
	defer resp.Body.Close()
	return
}

func PushPurchaseCompleted(creator *model.User, fan *model.User, payment *model.PendingPayment) {

}

func SendResetPasswordNotification(user model.User) {
	email := model.EmailNotification{
		Title:      "You have successfully reset your password",
		Template:   "account_change",
		ActionURL:  user.URL(),
		ActionName: "Go to your account",
		Email:      user.Email,
	}
	email.Owner = user.ID
	email.Message = "A change has been made to your account.\n\n\nYou changed your password through a password reset. If you did not do this, please contact support immediately."
	DispatchEmail(user, email, true)
	return
}

func SendPasswordChangedNotification(user model.User) {
	email := model.EmailNotification{
		Title:      "You changed your password",
		Template:   "account_change",
		ActionURL:  user.URL(),
		ActionName: "Go to your account",
		Email:      user.Email,
	}
	email.Owner = user.ID
	email.Message = "A change has been made to your account.\n\n\nYou changed your password. If you did not do this, please contact support immediately."
	DispatchEmail(user, email, true)
	return
}

func SendGettingStartedEmail(user model.User) (err error) {
	return
}

func SendSignupEmail(user model.User, session sessions.VisitorSession) (err error) {
	mail := model.EmailNotification{
		Email:      user.Email,
		Title:      "You have created an account with MyHustle",
		Template:   "account_signup",
		ActionURL:  user.URL(),
		ActionName: "View your account",
		Dictionary: session.Info.Dictionary(),
	}
	mail.Priority = true
	mail.Owner = user.ID
	DispatchEmail(user, mail, true)
	return
}

func SendPhoneNumberChangedEmail(user model.User) (err error) {
	email := model.EmailNotification{
		Title:      "Your phone number has been changed",
		Template:   "account_change",
		ActionURL:  user.URL(),
		ActionName: "Go to your account",
		Email:      user.Email,
	}
	email.Owner = user.ID
	email.Message = fmt.Sprintf("A change has been made to your account.\n\n\nYou changed your phone number to %s", user.PhoneNumber)
	DispatchEmail(user, email, true)
	return
}

func SendEmailChangedEmail(oldEmail string, user model.User) (err error) {
	email := model.EmailNotification{
		Title:      "Your email address has been changed",
		Template:   "account_change",
		ActionURL:  user.URL(),
		ActionName: "Go to your account",
		Email:      oldEmail,
	}
	email.Owner = user.ID
	email.Message = fmt.Sprintf("A change has been made to your account.\n\n\nYou changed your email address to %s", user.Email)
	DispatchEmail(user, email, true)
	return
}
