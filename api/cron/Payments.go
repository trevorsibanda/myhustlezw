package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/payments"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	_paymentsPendingQueue                    = "payments:dispatch"
	paymentsPendingQueue       *amqp.Queue   = nil
	paymentsPendingChannel     *amqp.Channel = nil
	paymentsScheduleCheckQueue               = "payments:upstream:check"
	paymentsUpdatedQueue                     = "payments:upstream:updates" //updates from upstream
)

type handleEcocashResponse struct {
	Status       string `json:"status"`
	Success      bool   `json:"success"`
	HasRedirect  bool   `json:"hasRedirect"`
	PollURL      string `json:"pollUrl"`
	Instructions string `json:"instructions"`
}

//DispatchPayment dispatches a payment in a seeprate goroutine
func ForceDispatchPaymentLocal(payment model.PendingPayment, session sessions.VisitorSession, creator model.User, node string) {
	//perform the transaction in the background
	var result []byte
	var err error
	result, err = sendPaymentRequest(payment)
	if err == nil {
		switch payment.Gateway {
		case model.GatewayEcocash:
			{
				var eres handleEcocashResponse
				//schedule a check using the pollURL
				if err = json.Unmarshal(result, &eres); err != nil || eres.Status != "ok" || !eres.Success {
					payment.Status = "cancelled"
					SendChannelMessage(payment.ID.Hex(), "cancelled", util.ScrubPublic(payment))
					_ = payment.UpdateStatus("cancelled", fmt.Sprintf("on node %s\n\n\n%s", node, string(result)))
					log.Printf("Transaction failed on node %s. Payment: %v Reason: %v", node, payment, string(result))
					return
				}
				payment.Status = "sent"
				payment.PollURL = eres.PollURL
				SendChannelMessage(payment.ID.Hex(), "sent", util.ScrubPublic(payment))
				_ = payment.UpdateStatus("sent", string(result))

				s := gocron.NewScheduler(time.UTC)
				paymentID := payment.ID
				job, _ := s.Every(30).Seconds().Do(func() {
					payment, err := model.RetrievePendingPayment(paymentID)
					status, logm, err1 := QueryStatus(payment)
					if err1 != nil {
						log.Printf("Failed to query status of payment %v %v %v %v", err1, status, logm, payment)
					} else if payment.Status != status {
						if status == "paid" {
							payment.Status = "paid"
							var support model.CreatorSupport
							if support, err = HandlePaymentSuccess(session, creator, payment); err != nil {
								fmt.Printf("Failed to handle payment success with error %v %v %v %v %v", session, creator, payment, err, support)
							}
							s.Stop()
						}
						err = payment.UpdateStatus(status, logm)
						payment.Status = status
						SendChannelMessage(payment.ID.Hex(), status, util.ScrubPublic(payment))
						return
					}
				})
				job.LimitRunsTo(6)
				s.StartAsync()

			}
		case model.Gateway2Checkout:
			{

			}
		default:

		}

	} else {
		payment.Status = "cancelled"
		SendChannelMessage(payment.ID.Hex(), "cancelled", util.ScrubPublic(payment))
		err = payment.UpdateStatus("cancelled", string(result))
		if err != nil {
			log.Printf("Failed to update payment status to cancelled %v, %v", err, payment)
		}
		return
	}
}

//QueryStatus checks the status of a pending payment
func QueryStatus(p model.PendingPayment) (status string, log string, err error) {
	var data []byte

	var resp payments.EcocashPollResponse
	var httpResponse *http.Response

	if p.PollURL != "" {
		payload, _ := json.Marshal(p)
		httpResponse, err = http.Post(p.PollURL, "application/x-www-form-urlencoded", bytes.NewBuffer(payload))

		if err == nil {
			defer httpResponse.Body.Close()
			data, _ = ioutil.ReadAll(httpResponse.Body)
		}

		if err = resp.FromFormData(string(data)); err != nil {
			return
		}

	} else {
		uri := payments.PaymentURI("poll", p)
		data, err = payments.ApiRequest(uri, p)
		if err = json.Unmarshal(data, &resp); err != nil {
			return
		}
	}

	log = string(data)
	if err != nil {
		return
	}

	switch p.Gateway {
	case model.GatewayEcocash:
		{
			status = strings.ToLower(resp.Status)
		}
	case model.Gateway2Checkout:
		{

		}
	}
	return
}

//HandlePaymentSuccess handles a payment when you have paid
//this method is called when payment status moves from pending/sent -> paid
func HandlePaymentSuccess(session sessions.VisitorSession, creator model.User, payment model.PendingPayment) (support model.CreatorSupport, err error) {
	var op model.WalletOperation
	if err = payment.UpdateStatus("paid", fmt.Sprintf("Session: %v\nCreator: %v\nPayment: %v", session, creator, payment)); err != nil {
		log.Printf("Failed to update payment status to paid. Error: %v  Session: %v\nCreator: %v\nPayment: %v", err, session, creator, payment)
	}
	payment.Status = "paid"

	if payment.Action != model.VerifyAccountOnPaid {
		if op, err = creator.CreditEscrow(payment.Price, string(payment.Gateway), string(payment.Currency), payment.Comment, session.User.ID); err != nil {
			//this is a critical failure
			log.Printf("WARNING: PAYMENTS: FAILED TO DEBIT USER: %v %v %v %v", session, creator, payment, err)
			//todo: notify user of error
			go SendDebitAccountFailedEmail(creator, payment)
			return
		}
	}

	supporter := model.AnonymousUser(payment.Email, fmt.Sprintf("%s%s", payment.PhoneCountryCode, payment.PhoneNationalNumber), payment.Fullname)
	supporter.ID = session.User.ID

	switch payment.Action {
	case model.AddSupportMessageOnPaid:
		{

			if support, err = creator.AddSupportMessage(supporter, payment.Comment, op, payment.Items, payment.ItemName, true); err == nil {
				//send confirmation emails and SMS
				err = OnAddSupportMessageOnPaid(supporter, creator, payment, support)
			} else {
				//todo : retry this
				log.Printf("Failed to add support message %v", err)
			}
		}
	case model.AllowServiceRequestOnPaid:
		{
			service, _ := model.CampaignByID(payment.TargetCampaign)
			service.ServiceDecrementQuantity()

			filledForm := payment.Form //campaign.Form.Fill(payment.Form.Email, payment.Form.Phone, supporter.Fullname, payment.Form.Answer, payment.Form.SelectedOption, payment.Form.ExtraInfo, false)
			if support, err = creator.AddServiceRequest(supporter, service, filledForm, payment, op); err == nil {
				payment.UnlockCode = service.URL(creator.Username, &support.ID, &support.UnlockCode)
				err = OnAddServiceRequest(supporter, creator, service, payment, support)
			}
		}
	case model.AllowAccessOnPaid:
		{
			//add session to redis

			campaign, _ := model.CampaignByID(payment.TargetCampaign)

			if support, err = creator.AddPayPerViewAccess(campaign, supporter, op); err == nil {
				payment.UnlockCode = campaign.URL(creator.Username, &support.ID, &support.UnlockCode)
				err = OnPayPerViewAccessGranted(supporter, campaign, creator, payment, support)
			} else {
				//todo : retry this
				log.Printf("Failed to add access request %v %v %v %v", err, payment, campaign, supporter)
			}
		}
	case model.SubscribeOnPaid:
		{

			if support, err = creator.AddSubscriber(supporter, op, payment.TargetCampaign, !payment.Anonymous); err == nil {
				//send confirmation emails and SMS
				go SendNewSubscriberNotification(supporter, creator, payment, op, support)
			}

		}
	case model.VerifyAccountOnPaid:
		{
			if err = creator.VerifyIdentity(); err == nil {
				go SendIdentityVerifiedNotification(&creator)
			} else {
				log.Printf("Failed to verify identity %v for %v", err, creator)
			}
		}
	default:
		err = fmt.Errorf("Unkown payment action on payment %v %v", payment.Action, payment)
	}

	SendChannelMessage(payment.ID.Hex(), "paid", util.ScrubPublic(payment))
	return
}

type pendingPaymentDispatch struct {
	Payment model.PendingPayment
	Session sessions.VisitorSession
	User    model.User
	Node    string
}

//DispatchPayment dispatches a payment request
func DispatchPayment(p model.PendingPayment, s sessions.VisitorSession, u model.User) {
	entry := pendingPaymentDispatch{
		Payment: p,
		Session: s,
		User:    u,
		Node:    runnerNode,
	}

	if bsonEntry, err := bson.Marshal(entry); err != nil {
		log.Printf("payments Failed to bson marshal entry %v", err)
	} else {
		if err := paymentsPendingChannel.Publish("", paymentsPendingQueue.Name, false, false, amqp.Publishing{
			ContentType:  "application/bson",
			Body:         bsonEntry,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			//send support email
			SendEngineerSupportEmail(fmt.Sprintf("Failed to publish pending payment to queue %v . Force initiating! ", err), p)
			ForceDispatchPaymentLocal(p, s, u, runnerNode)
		}
	}
}

func sendPaymentRequest(p model.PendingPayment) (res []byte, err error) {
	uri := payments.PaymentURI("dispatch", p)
	res, err = payments.ApiRequest(uri, p)
	return
}

//PaymentsDispatcher looks for queued payments and creates them
func paymentsDispatcherRunner() {
	log.Printf("payments dispatcher started at %s", time.Now())
	assert(paymentsPendingQueue != nil, "paymentsPending queue is nil")

	defer paymentsPendingChannel.Close()

	msgs, err := paymentsPendingChannel.Consume(
		_paymentsPendingQueue, // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		panic(fmt.Sprintf("paymentsPending Failed to consume from queue %v", err))
	}
	for d := range msgs {
		var entry pendingPaymentDispatch
		if err := bson.Unmarshal([]byte(d.Body), &entry); err != nil {
			log.Printf("payments Failed to unmarshal message %v", err)
			break
		}
		log.Printf("payments Processing payment %v", entry.Payment.ID)
		ForceDispatchPaymentLocal(entry.Payment, entry.Session, entry.User, entry.Node)
		d.Ack(false)
		log.Printf("payments Waiting for payment to process")

	}
}

//PaymentsSettlerRunner rnuns at predefined intervals looking for payments which can be settled
func PaymentsSettlerRunner() {
	var escrowPayments []model.WalletOperation
	var creators map[primitive.ObjectID]model.User = make(map[primitive.ObjectID]model.User)
	var err error

	log.Printf("paymentsettler running. Started at %s", time.Now())

	escrowPayments, err = model.RetrieveEscrowReadyWalletOps(1000)
	if err != nil {
		log.Printf("paymentsettler Failed to get escrow payments %v", err)
		return
	}

	for _, payment := range escrowPayments {
		var creator model.User
		var err error
		if creator, err = model.RetrieveCreatorByID(payment.Creator); err != nil {
			log.Printf("paymentsettler failed to retrieve creator for paymnet %v", payment)
			continue
		}

		err = payment.Available(fmt.Sprintf("\nOP: %v\nCreator: %v\nTime: %v", payment, creator, time.Now()))
		if err != nil {
			log.Printf("paymentsettler Failed to mark a payment as now available %v ERROR: %v", payment, err)
		} else {
			creators[payment.Creator] = creator
		}
	}

	for _, creator := range creators {
		usdSummary, _ := creator.WalletSummary(string(model.USD))
		zwlSummary, _ := creator.WalletSummary(string(model.ZWL))

		message := fmt.Sprintf(
			"You have %s and %s that is ready to be withdrawn from your wallet. ",
			FormatAsUSD(usdSummary.Available),
			FormatAsZWL(zwlSummary.Available))

		if usdSummary.Escrow > 0 || zwlSummary.Escrow > 0 {
			message = fmt.Sprintf("%s\nYou will be notified when the %s and %s in escrow becomes available to withdraw", message,
				FormatAsUSD(usdSummary.Escrow),
				FormatAsZWL(zwlSummary.Escrow))
		}

		if usdSummary.PendingWithdrawal > 0 || zwlSummary.PendingWithdrawal > 0 {
			message = fmt.Sprintf("%s.\n\n You also have %s and %s withdrawals that are been processed", message,
				FormatAsUSD(usdSummary.PendingWithdrawal),
				FormatAsZWL(zwlSummary.PendingWithdrawal))
		}

		if usdSummary.Disputed > 0 || zwlSummary.Disputed > 0 {
			message = fmt.Sprintf("%s\nThe %s and %s in disputed funds will be automatically released to you or refunded once resolved",
				message,
				FormatAsUSD(usdSummary.Disputed),
				FormatAsZWL(zwlSummary.Disputed))
		}

		SendWalletChangeEmail("You have funds ready to withdraw.", message, creator)
	}

	log.Printf("paymentsettler completed. Stopped at %s. Processed %d payments", time.Now(), len(escrowPayments))
	return
}

//PaymentsDisputeHandler runs at predefined intervals looking for disputed payments to report on
func PaymentsDisputeHandler() {

}

//PaymentsCashoutReminder runs at predefined intervals looking for accounts which can be cashed out
func PaymentsCashoutReminder() {

}

//PaymentsDailyNotifications runs once daily and notifies creators who made earnings on a particular day of
// their earning
func PaymentsDailyNotifications() {

}

//PaymentsDailyAdministrativeReport sends an administrative report outline of all payments received in a day
func PaymentsDailyAdministrativeReport() {

}

//PaymentsPendingGarabageCollector runs at predefined intervals and archives any pending payments which have not been fulfilled
func PaymentsPendingGarabageCollector() {

}
