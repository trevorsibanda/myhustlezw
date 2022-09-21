package api

import (
	"crypto/md5"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/controllers"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"github.com/ttacon/libphonenumber"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type initiateSupportForm struct {
	Gateway           string `form:"gateway" binding:"required"`
	Creator           string `form:"creator" binding:"required"`
	Items             int    `form:"items" binding:"required"` //number of items
	Supporter         string `form:"supporter"`
	Campaign          string `form:"campaign"`
	Phone             string `form:"phone" binding:"required"`
	NotificationPhone string `form:"notification_phone"`
	Fullname          string `form:"fullname"`
	Email             string `form:"email" binding:"required"`
	Message           string `form:"message" `
	Password          string `form:"password"`
}

//InitiateSupportPayment initiates a support/donation payment
func InitiateSupportPayment(ctx *gin.Context) {
	purpose, _ := ctx.Params.Get("purpose")
	switch purpose {
	case "support":
		InitiateDonationPayment(ctx)
	case "service":
		InitiateServicePayment(ctx)
	case "pay_per_view":
		InitiatePayPerViewPayment(ctx)
	case "subscribe":
		InitiateSubscriptionPayment(ctx)
	default:
		apiError(ctx, fmt.Sprintf("unknown payment type %s", purpose))
	}
	return
}

//InitiateDonationPayment handles a donation payment
func InitiateDonationPayment(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)

	if err != nil {
		session = sessions.NewVisitorSession(ctx)
	}

	var form initiateSupportForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid form %s", err))
		return
	}

	log.Printf("%v", form)

	if form.Items <= 0 {
		apiError(ctx, fmt.Sprintf("Need at least one item"))
		return
	}

	form.Email = strings.TrimSpace(strings.ToLower(form.Email))
	if !isEmailValid(form.Email) {
		apiError(ctx, fmt.Sprintf("%s is not a valid email", form.Email))
		return
	}

	var creator model.User
	var supporter model.User

	if creator, err = model.RetrieveCreatorByUsername(form.Creator); err != nil {
		apiError(ctx, fmt.Sprintf("User %s does not exist %s", form.Creator, err))
		return
	}

	if !creator.Page.AllowSupporters {
		apiError(ctx, "This user is not accepting support in this form.")
		return
	}

	if supporterID, err1 := primitive.ObjectIDFromHex(form.Supporter); err1 != nil {
		supporter = model.AnonymousUser(form.Email, form.Phone, form.Fullname)
	} else {
		if supporterID.IsZero() {
			supporter = model.AnonymousUser(form.Email, form.Phone, form.Fullname)
		} else if supporter, err = model.RetrieveCreatorByID(supporterID); err != nil {
			apiError(ctx, fmt.Sprintf("Supporter %s does not exist", form.Supporter))
			return
		}
	}

	var campaignID primitive.ObjectID
	if campaignID, err = primitive.ObjectIDFromHex(form.Campaign); err != nil {
		campaignID = primitive.NilObjectID
	} else {
		if c, err := model.CampaignByID(campaignID); err != nil {
			apiError(ctx, fmt.Sprintf("Campaign %s does not exist", form.Campaign))
			return
		} else if c.Owner != creator.ID {
			apiError(ctx, fmt.Sprintf("Campaign %s does not belong to user %s", form.Campaign, form.Creator))
			return
		}
	}

	var gateway model.PaymentGateway
	var currency model.PaymentCurrency
	var price float64
	var expires time.Time
	switch form.Gateway {
	case "ecocash":
		gateway = model.GatewayEcocash
		currency = model.ZWL
		price = cron.USDToZWL(creator.Page.DonationItemUnitPrice * float64(form.Items))
		expires = time.Now().Add(time.Minute * 10)
	case "international":
		gateway = model.Gateway2Checkout
		currency = model.USD
		price = creator.Page.DonationItemUnitPrice
		expires = time.Now().Add(time.Minute * 30)
	default:
		apiError(ctx, fmt.Sprintf("Invalid payment gateway %v", form.Gateway))
		return
	}

	//parse phone number
	if gateway == model.GatewayEcocash {
		if strings.HasPrefix(form.Phone, "+") {
			if !strings.HasPrefix(form.Phone, "+263") {
				apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", form.Phone))
				return
			}
		}
		if strings.HasPrefix(form.Phone, "0") {
			if strings.HasPrefix(form.Phone, "0263") || strings.HasPrefix(form.Phone, "00263") {
				apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", form.Phone))
				return
			}
			if !strings.HasPrefix(form.Phone, "078") && !strings.HasPrefix(form.Phone, "077") {
				apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", form.Phone))
				return
			}
			form.Phone = "+263" + form.Phone[1:]
		}
	}

	var phone *libphonenumber.PhoneNumber
	if phone, err = libphonenumber.Parse(form.Phone, "ZW"); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid phone number passed"))
		return
	}

	if gateway == model.GatewayEcocash && phone.GetCountryCode() != 263 {
		apiError(ctx, fmt.Sprintf("Only zimbabwean phone numbers allowed for mobile money payment"))
		return
	}

	if number := fmt.Sprintf("%d", phone.GetNationalNumber()); gateway == model.GatewayEcocash && !strings.HasPrefix(number, "77") && !strings.HasPrefix(number, "78") {
		apiError(ctx, fmt.Sprintf("Only Econet phone numbers can be used for Ecocash payments"))
		return
	}

	form.Fullname = strings.TrimSpace(form.Fullname)
	if len(form.Fullname) > 150 {
		form.Fullname = form.Fullname[0:150]
	}

	var paymentID primitive.ObjectID
	payment := model.PendingPayment{
		Anonymous:           supporter.ID.IsZero(),
		Created:             time.Now(),
		Comment:             strings.TrimSpace(form.Message),
		Gateway:             gateway,
		Email:               supporter.Email,
		Fullname:            supporter.Fullname,
		PhoneCountryCode:    fmt.Sprintf("+%d", phone.GetCountryCode()),
		PhoneNationalNumber: fmt.Sprintf("%d", phone.GetNationalNumber()),
		ExtraPhone:          form.NotificationPhone,
		Currency:            currency,
		Items:               form.Items,
		ItemName:            creator.Page.DonationItemName,
		Price:               price,
		TargetCreator:       creator.ID,
		TargetFan:           supporter.ID,
		TargetCampaign:      campaignID,
		ThankYou:            creator.Page.ThankYouMessage,
		Action:              model.AddSupportMessageOnPaid,
		UpsertFan:           false,
		Status:              "queued",
	}

	if paymentID, err = model.NewDonationPendingPayment(payment); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save new payment"))
		return
	}
	payment.ID = paymentID

	md5Sum := md5.Sum([]byte(fmt.Sprintf("%snonce%d%s%d", payment.ID.Hex(), expires.Unix(), payment.ID.Hex(), expires.Unix())))

	payment.ThankYou = ""
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"_id":       paymentID,
		"payment":   payment,
		"ts":        expires.Unix(),
		"signature": fmt.Sprintf("%x", md5Sum),
	}, false)

	go cron.DispatchPayment(payment, session, creator)
}

type paynowCallbackForm struct {
	Status         string  `form:"status"`
	Reference      string  `form:"reference"`
	PaynowRefernce string  `form:"paynowreference"`
	PollURL        string  `form:"pollurl"`
	Amount         float32 `form:"amount"`
	Hash           string  `form:"hash"`
}

//EcocashTransactionCallback handles callback events from paynow/ecocash
func EcocashTransactionCallback(ctx *gin.Context) {
	var form paynowCallbackForm
	var paymentID primitive.ObjectID
	var creator model.User

	session := sessions.NewVisitorSession(ctx)

	var err error
	if err := ctx.Bind(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to parse callback form"))
		return
	}

	//paymentID
	if paymentID, err = primitive.ObjectIDFromHex(form.Reference); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid payment ID"))
		return
	}

	//retrieve payment
	var payment model.PendingPayment
	if payment, err = model.RetrievePendingPayment(paymentID); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to retrieve payment"))
		return
	}

	//load creator
	if creator, err = model.RetrieveCreatorByID(payment.TargetCreator); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to retrieve creator"))
		return
	}

	//load user
	var user model.User
	if user, err = model.RetrieveCreatorByID(payment.TargetFan); err != nil {
		user = model.AnonymousUser(payment.Email, payment.ExtraPhone, payment.Fullname)
	}
	session.User = user

	//check hash
	//TODO

	status := strings.ToLower(form.Status)
	log.Printf("Ecocash callback. updated payment %v to %s with form %v", paymentID, status, form)
	if status == "paid" && payment.Status != "paid" {
		payment.Status = "paid"
		go cron.HandlePaymentSuccess(session, creator, payment)
		util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "updated_paid"}, false)
		return
	}

	if status == "cancelled" && payment.Status != "cancelled" {
		cron.SendChannelMessage(paymentID.Hex(), "cancelled", util.ScrubPublic(payment))
		payment.UpdateStatus("cancelled", fmt.Sprintf("%v", form))
	}

	if status == "failed" && payment.Status != "cancelled" {
		cron.SendChannelMessage(paymentID.Hex(), "cancelled", util.ScrubPublic(payment))
		payment.UpdateStatus("cancelled", fmt.Sprintf("%v", form))
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "updated_nothing"}, false)

}

//InitiateSubscriptionPayment initiates a payment towards a user's subscriptions
func InitiateSubscriptionPayment(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)

	if err != nil {
		session = sessions.NewVisitorSession(ctx)
	}

	var form initiateSupportForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid form %s", err))
		return
	}

	//upsert user if does not exist
	if session.User.ID.IsZero() {
		//todo upsert fan
		session.User.Email = form.Email
		session.User.PhoneNumber = form.Phone
		session.User.Fullname = form.Fullname

		if user, err := model.RetrieveCreator(form.Email, form.Phone, ""); err == nil {
			//validate credentials if matching login otherwise err
			sec, err := model.RetrieveUserCredentials(user.ID)
			if err == nil && model.ValidatePasswords(form.Password, sec.Password) {
				session.User = user
				session.HasCreator = true
				session.LastActive = time.Now()
				session.User.LoggedIn = true
				session.Save(ctx)
				go cron.PushNotifyNewCreatorLogin(user, session, &sec)
			} else {
				apiError(ctx, "An account with that email or phone number already exists. Please login")
				return
			}
		} else {
			//create new user
			//check if a user exists, otherwise signup
			session.User.PhoneVerified = true
			session.User.EmailVerified = true
			session.User.LoggedIn = true
			session.HasCreator = true
			session.LastActive = time.Now()
			id := primitive.NewObjectID()
			s := fmt.Sprintf("%x", md5.Sum(id[:]))
			user, _, err := model.NewCreator(session.User, s, form.Password)
			if err != nil {
				apiError(ctx, fmt.Sprintf("Failed to create new user with error %v", err))
				return
			} else {
				go cron.SendSignupEmail(*user, session)
				session.Save(ctx)
			}
		}
	}

	var gateway model.PaymentGateway
	var currency model.PaymentCurrency
	switch form.Gateway {
	case "ecocash":
		gateway = model.GatewayEcocash
		currency = model.ZWL
	case "international":
		gateway = model.Gateway2Checkout
		currency = model.USD
	default:
		apiError(ctx, fmt.Sprintf("Invalid payment gateway %v", form.Gateway))
		return
	}

	var creator model.User

	if creator, err = model.RetrieveCreatorByUsername(form.Creator); err != nil {
		apiError(ctx, fmt.Sprintf("User %s does not exist %s", form.Creator, err))
		return
	}

	campaignID, _ := primitive.ObjectIDFromHex(form.Campaign)

	//parse phone number
	var phone *libphonenumber.PhoneNumber
	if phone, err = libphonenumber.Parse(form.Phone, "ZW"); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid phone number passed"))
		return
	}

	if gateway == model.GatewayEcocash && phone.GetCountryCode() != 263 {
		apiError(ctx, fmt.Sprintf("Only zimbabwean phone numbers allowed for mobile money payment"))
		return
	}

	if number := fmt.Sprintf("%d", phone.GetNationalNumber()); gateway == model.GatewayEcocash && !strings.HasPrefix(number, "77") && !strings.HasPrefix(number, "78") {
		apiError(ctx, fmt.Sprintf("Only Econet phone numbers can be used for Ecocash payments"))
		return
	}

	var paymentID primitive.ObjectID
	payment := model.PendingPayment{
		Created:             time.Now(),
		Comment:             form.Message,
		Gateway:             gateway,
		Email:               session.User.Email,
		Fullname:            session.User.Fullname,
		PhoneCountryCode:    fmt.Sprintf("+%d", phone.GetCountryCode()),
		PhoneNationalNumber: fmt.Sprintf("%d", phone.GetNationalNumber()),
		Currency:            currency,
		Items:               creator.Subscriptions.Period,
		ItemName:            "Subscription",
		Price:               cron.USDToZWL(creator.Subscriptions.Price),
		TargetCreator:       creator.ID,
		TargetFan:           session.User.ID,
		TargetCampaign:      campaignID,
		ThankYou:            "",
		Action:              model.SubscribeOnPaid,
		UpsertFan:           true,
		Status:              "queued",
		State:               form.Password,
	}

	//don't double subscribe
	if canView, sub := controllers.CanViewCreatorSubscriberContent(creator, session, nil, false); canView {
		if sub.ID.IsZero() {
			apiError(ctx, "You can view this content but you are not subscribed")
			return
		} else {
			apiError(ctx, "You are already subscribed")
			return
		}

	}
	if paymentID, err = model.NewSubscriptionPendingPayment(payment); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save new payment"))
		return
	}
	payment.ID = paymentID

	expires := time.Now().Add(time.Minute * 10)
	md5Sum := md5.Sum([]byte(fmt.Sprintf("%snonce%d%s%d", payment.ID.Hex(), expires.Unix(), payment.ID.Hex(), expires.Unix())))

	go cron.DispatchPayment(payment, session, creator)

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"_id":       paymentID,
		"payment":   payment,
		"ts":        expires.Unix(),
		"signature": fmt.Sprintf("%x", md5Sum),
	}, false)
}

type initiatePayPerViewForm struct {
	Gateway           string `form:"gateway" binding:"required"`
	Creator           string `form:"creator" binding:"required"`
	Supporter         string `form:"supporter" binding:"required"`
	Campaign          string `form:"campaign" binding:"required"`
	Phone             string `form:"phone" binding:"required"`
	NotificationPhone string `form:"notification_phone"`
	Fullname          string `form:"fullname"`
	Email             string `form:"email" binding:"required"`
}

//InitiatePayPerViewPayment initiates a pay per view payment
func InitiatePayPerViewPayment(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)

	if err != nil {
		session = sessions.NewVisitorSession(ctx)
	}

	var form initiatePayPerViewForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid form %s", err))
		return
	}

	form.Email = strings.TrimSpace(strings.ToLower(form.Email))
	if !isEmailValid(form.Email) {
		apiError(ctx, fmt.Sprintf("%s is not a valid email", form.Email))
		return
	}

	var creator model.User
	var supporter model.User = session.User

	if creator, err = model.RetrieveCreatorByUsername(form.Creator); err != nil {
		apiError(ctx, fmt.Sprintf("User %s does not exist %s", form.Creator, err))
		return
	}

	var campaignID primitive.ObjectID
	var c model.Campaign
	if campaignID, err = primitive.ObjectIDFromHex(form.Campaign); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid campaign ID %s", form.Campaign))
		return
	} else {
		if c, err = model.CampaignByID(campaignID); err != nil {
			apiError(ctx, fmt.Sprintf("Content %s does not exist", form.Campaign))
			return
		} else if c.Owner != creator.ID {
			apiError(ctx, fmt.Sprintf("Content %s does not belong to user %s", form.Campaign, form.Creator))
			return
		} else if c.Subscription != "pay_per_view" {
			apiError(ctx, fmt.Sprintf("Content is not available as pay per view anymore"))
			return
		} else if !c.Active {
			apiError(ctx, fmt.Sprintf("This content has been deleted"))
			return
		} else if c.ExpiresAt.After(time.Now()) {
			apiError(ctx, fmt.Sprintf("This content has expired"))
			return
		}
	}

	if !session.User.ID.IsZero() {
		if supports, _ := c.FetchSupporters(session.User.ID, model.ServiceRequest); len(supports) != 0 {
			apiError(ctx, "You have already paid to access this content.")
			return
		}
	}

	var gateway model.PaymentGateway
	var currency model.PaymentCurrency
	var price float64
	var expires time.Time
	switch form.Gateway {
	case "ecocash":
		gateway = model.GatewayEcocash
		currency = model.ZWL
		price = cron.USDToZWL(c.Price)
		expires = time.Now().Add(time.Minute * 10)
	case "international":
		gateway = model.Gateway2Checkout
		currency = model.USD
		price = c.Price
		expires = time.Now().Add(time.Minute * 30)
	default:
		apiError(ctx, fmt.Sprintf("Invalid payment gateway %v", form.Gateway))
		return
	}

	//parse phone number
	if gateway == model.GatewayEcocash {
		if strings.HasPrefix(form.Phone, "+") && !strings.HasPrefix(form.Phone, "+263") {
			apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", form.Phone))
			return

		}
		if strings.HasPrefix(form.Phone, "0") {
			if strings.HasPrefix(form.Phone, "0263") || strings.HasPrefix(form.Phone, "00263") {
				apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", form.Phone))
				return
			}
			if !strings.HasPrefix(form.Phone, "078") && !strings.HasPrefix(form.Phone, "077") {
				apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", form.Phone))
				return
			}
			form.Phone = "+263" + form.Phone[1:]
		}
	}

	var phone *libphonenumber.PhoneNumber
	if phone, err = libphonenumber.Parse(form.Phone, "ZW"); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid phone number passed"))
		return
	}

	if gateway == model.GatewayEcocash && phone.GetCountryCode() != 263 {
		apiError(ctx, fmt.Sprintf("Only Zimbabwean phone numbers allowed for mobile money payment"))
		return
	}

	if number := fmt.Sprintf("%d", phone.GetNationalNumber()); gateway == model.GatewayEcocash && !strings.HasPrefix(number, "77") && !strings.HasPrefix(number, "78") {
		apiError(ctx, fmt.Sprintf("Only Econet phone numbers can be used for Ecocash payments"))
		return
	}

	form.Fullname = strings.TrimSpace(form.Fullname)
	if len(form.Fullname) > 150 {
		form.Fullname = form.Fullname[0:150]
	}
	if supporter.ID.IsZero() {
		supporter.Fullname = form.Fullname
		supporter.Email = form.Email
	}

	var paymentID primitive.ObjectID
	payment := model.PendingPayment{
		Anonymous:           supporter.ID.IsZero(),
		Created:             time.Now(),
		Gateway:             gateway,
		Email:               supporter.Email,
		Fullname:            supporter.Fullname,
		PhoneCountryCode:    fmt.Sprintf("+%d", phone.GetCountryCode()),
		PhoneNationalNumber: fmt.Sprintf("%d", phone.GetNationalNumber()),
		ExtraPhone:          form.NotificationPhone,
		Currency:            currency,
		Price:               price,
		TargetCreator:       creator.ID,
		TargetFan:           supporter.ID,
		TargetCampaign:      campaignID,
		ItemName:            c.Title,
		Action:              model.AllowAccessOnPaid,
		UpsertFan:           false,
		Status:              "queued",
	}

	if paymentID, err = model.NewPayPerViewPendingPayment(payment); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save new payment"))
		return
	}
	payment.ID = paymentID

	md5Sum := md5.Sum([]byte(fmt.Sprintf("%snonce%d%s%d", payment.ID.Hex(), expires.Unix(), payment.ID.Hex(), expires.Unix())))

	payment.ThankYou = ""
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"_id":       paymentID,
		"payment":   payment,
		"ts":        expires.Unix(),
		"signature": fmt.Sprintf("%x", md5Sum),
	}, false)

	go cron.DispatchPayment(payment, session, creator)
}

type initiateServicePaymentForm struct {
	Gateway      string `form:"gateway" binding:"required"`
	Creator      string `form:"creator" binding:"required"`
	Service      string `form:"service" binding:"required"`
	Phone        string `form:"phone" binding:"required"`
	Ecocash      string `form:"ecocash" binding:"required"`
	Fullname     string `form:"fullname" binding:"required"`
	Email        string `form:"email" binding:"required"`
	Question     string `form:"question"`
	OptionChoice string `form:"answer" json:"answer" binding:"required"`
	ExtraInfo    string `form:"extra_info" `
}

//InitiateServicePayment initiates a payment for a service
func InitiateServicePayment(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)

	if err != nil {
		session = sessions.NewVisitorSession(ctx)
	}

	var form initiateServicePaymentForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid form %s", err))
		return
	}

	var campaign model.Campaign
	var campaignID primitive.ObjectID
	if campaignID, err = primitive.ObjectIDFromHex(form.Service); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid campaign id passed %v", form.Service))
		return
	}
	if campaign, err = model.CampaignByID(campaignID); err != nil {
		apiError(ctx, fmt.Sprintf("Campaign does not exist %v", campaignID))
		return
	}

	if campaign.Type != model.ServiceCampaign {
		apiError(ctx, fmt.Sprintf("Campaign %v is not a service", campaignID))
		return
	}

	if campaign.Form.QuantityAvailable == 0 {
		apiError(ctx, "Sorry there are no more items available to create this order :(")
		return
	}

	var creator model.User

	if creator, err = model.RetrieveCreatorByUsername(form.Creator); err != nil {
		apiError(ctx, fmt.Sprintf("User %s does not exist %s", form.Creator, err))
		return
	}

	if creator.ID != campaign.Owner {
		apiError(ctx, fmt.Sprintf("User %s does not own campaign %v", form.Creator, campaignID))
		return
	}

	var gateway model.PaymentGateway
	var currency model.PaymentCurrency
	var price float64
	var expires time.Time
	switch form.Gateway {
	case "ecocash":
		gateway = model.GatewayEcocash
		currency = model.ZWL
		price = cron.USDToZWL(campaign.Price)

		expires = time.Now().Add(time.Minute * 10)
	case "international":
		gateway = model.Gateway2Checkout
		currency = model.USD
		price = campaign.Price
		expires = time.Now().Add(time.Minute * 30)
	default:
		apiError(ctx, fmt.Sprintf("Invalid payment gateway %v", form.Gateway))
		return
	}

	form.Email = strings.TrimSpace(strings.ToLower(form.Email))
	if !isEmailValid(form.Email) {
		apiError(ctx, fmt.Sprintf("%s is not a valid email", form.Email))
		return
	}

	//parse phone number
	ecocashNumber := form.Ecocash

	if ecocashNumber == "" {
		ecocashNumber = form.Phone
	}

	var phone *libphonenumber.PhoneNumber
	if phone, err = libphonenumber.Parse(ecocashNumber, "ZW"); err != nil {
		apiError(ctx, "Invalid phone number passed")
		return
	}

	if gateway == model.GatewayEcocash && phone.GetCountryCode() != 263 {
		apiError(ctx, "Only zimbabwean phone numbers allowed for mobile money payment")
		return
	}

	if number := fmt.Sprintf("%d", phone.GetNationalNumber()); gateway == model.GatewayEcocash && !strings.HasPrefix(number, "77") && !strings.HasPrefix(number, "78") {
		apiError(ctx, "Only Econet phone numbers can be used for Ecocash payments")
		return
	}

	form.Fullname = strings.TrimSpace(form.Fullname)
	if len(form.Fullname) > 150 {
		form.Fullname = form.Fullname[0:150]
	}

	serviceForm := campaign.Form.Fill(form.Email,
		fmt.Sprintf("+%d%d", phone.GetCountryCode(), phone.GetNationalNumber()),
		form.Fullname,
		form.OptionChoice,
		"none",
		form.ExtraInfo,
		false)
	var paymentID primitive.ObjectID
	payment := model.PendingPayment{
		Created:             time.Now(),
		Gateway:             gateway,
		Email:               form.Email,
		Fullname:            form.Fullname,
		PhoneCountryCode:    fmt.Sprintf("+%d", phone.GetCountryCode()),
		PhoneNationalNumber: fmt.Sprintf("%d", phone.GetNationalNumber()),
		ExtraPhone:          form.Phone,
		Currency:            currency,
		Items:               1,
		ItemName:            campaign.Title,
		Price:               price,
		TargetCreator:       creator.ID,
		TargetCampaign:      campaign.ID,
		TargetFan:           session.User.ID,
		ThankYou:            campaign.Form.ThankYouMessage,
		Action:              model.AllowServiceRequestOnPaid,
		UpsertFan:           false,
		Form:                &serviceForm,
		Status:              "queued",
	}

	if paymentID, err = model.NewServicePendingPayment(payment); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save new payment"))
		return
	}
	payment.ID = paymentID

	md5Sum := md5.Sum([]byte(fmt.Sprintf("%snonce%d%s%d", payment.ID.Hex(), expires.Unix(), payment.ID.Hex(), expires.Unix())))

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"_id":       paymentID,
		"payment":   payment,
		"ts":        expires.Unix(),
		"signature": fmt.Sprintf("%x", md5Sum),
	}, false)

	cron.DispatchPayment(payment, session, creator)

}

//PollPaymentStatus polls a payment's status
func PollPaymentStatus(ctx *gin.Context) {

	session, err := sessions.GetVisitorSession(ctx)

	if err != nil {
		session = sessions.NewVisitorSession(ctx)
	}

	session.LastActiveNow()
	session.Save(ctx)

	var idString string
	var ts string
	var signature string
	var ok bool
	var id primitive.ObjectID

	if idString, ok = ctx.Params.Get("id"); !ok {
		apiError(ctx, fmt.Sprintf("No ID passed"))
		return
	} else if id, err = primitive.ObjectIDFromHex(idString); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid ID passed"))
		return
	}

	ts, _ = ctx.GetQuery("ts")
	if tsInt, _ := strconv.Atoi(ts); time.Now().Unix() > int64(tsInt) {
		apiError(ctx, fmt.Sprintf("payment url has expired. create a new payment"))
		return
	}

	signature, _ = ctx.GetQuery("signature")

	md5Sum := md5.Sum([]byte(fmt.Sprintf("%snonce%s%s%s", idString, ts, idString, ts)))
	if fmt.Sprintf("%x", md5Sum) != signature {
		apiError(ctx, fmt.Sprintf("Invalid signature. Do not tamper with URL params"))
		return
	}

	var payment model.PendingPayment

	if payment, err = model.RetrievePendingPayment(id); err != nil {
		apiError(ctx, fmt.Sprintf("Payment does not exist"))
		return
	}

	if payment.Status != "paid" {
		payment.ThankYou = ""
	}

	if payment.Status == "paid" && payment.Action == model.AllowAccessOnPaid {
		//update session to include new id
		target := payment.TargetCampaign
		session.UpdateGrantContentAccess(ctx, &target)
	}

	util.ScrubbedPublicAPIJSON(ctx, payment, false)

}

//HandleCallbackEcocash handles callback events for ecocash
func HandleCallbackEcocash(ctx *gin.Context) {

}

//HandleCallbackPaynow handles callback events for paynow
func HandleCallbackPaynow(ctx *gin.Context) {

}

//HandleCallback2Checkout handles IPN from 2checkout
func HandleCallback2Checkout(ctx *gin.Context) {

}
