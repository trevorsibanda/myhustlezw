package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"github.com/ttacon/libphonenumber"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//GetCurrentUser returns the current user
func GetCurrentUser(ctx *gin.Context) {
	var session sessions.VisitorSession
	var creator model.User
	var loggedIn bool
	var err error

	session, err = sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
		session.Save(ctx)

	}
	if creator, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		creator = model.AnonymousUser("anonymous@myhustle.co.zw", "", "Anonymous")
		loggedIn = false
	} else {
		loggedIn = true
	}

	creator.LoggedIn = loggedIn
	util.ScrubbedPublicAPIJSON(ctx, creator, loggedIn)
}

func UpdateUserWebSubscription(ctx *gin.Context) {
	session := ctx.Keys["session"].(sessions.VisitorSession)

	var data []byte
	var err error

	if data, err = ioutil.ReadAll(ctx.Request.Body); err != nil {
		apiError(ctx, "Failed to read request body")
		return
	}

	creator, _ := model.RetrieveCreatorByID(session.User.ID)

	if err = creator.SetPushSubscription(data, false); err != nil {
		apiError(ctx, "Failed to save subscription")
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "ok",
	}, true)

}

type basicPageDetailsForm struct {
	Username     string `form:"username" binding:"required"`
	Fullname     string `form:"fullname" binding:"required"`
	Description  string `form:"description" binding:"required"`
	Aboutme      string `form:"aboutme" binding:"required"`
	PersonalSite string `bson:"url" form:"url" json:"url" groups:"public"`
	Facebook     string `bson:"facebook" form:"facebook" json:"facebook" groups:"public"`
	Twitter      string `bson:"twitter" form:"twitter" json:"twitter" groups:"public"`
	Instagram    string `bson:"instagram" form:"instagram" json:"instagram" groups:"public"`
	Whatsapp     string `bson:"whatsapp" form:"whatsapp" json:"whatsapp" groups:"public"`
	Youtube      string `bson:"youtube" form:"youtube" json:"youtube" groups:"public"`
}

func validateSocialMedia(links *model.CreatorSocialMedia) {

	links.Instagram = strings.ToLower(links.Instagram)
	if len(links.Instagram) > 15 {
		links.Instagram = links.Instagram[:15]
	}

	links.Twitter = strings.TrimSpace(strings.ToLower(links.Twitter))
	if len(links.Twitter) > 20 {
		links.Twitter = links.Twitter[:20]
	}

	_, err := url.ParseRequestURI(links.PersonalSite)
	if err != nil {
		links.PersonalSite = ""
	}

	_, err = url.ParseRequestURI(links.Facebook)
	if err != nil {
		links.Facebook = ""
	}
	if !strings.Contains(links.Facebook, "facebook.com") && !strings.Contains(links.Facebook, "fb.com") {
		links.Facebook = ""
	}

	_, err = url.ParseRequestURI(links.Youtube)
	if err != nil {
		links.Youtube = ""
	}
	if !strings.Contains(links.Youtube, "youtube.com") && !strings.Contains(links.Facebook, "youtu.be") {
		links.Youtube = ""
	}

}

//UpdateBasicUserPageDetails does as the name explains
func UpdateBasicUserPageDetails(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	var form basicPageDetailsForm
	var err error
	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, err)
		return
	}

	l := len(form.Aboutme)
	if l > 150 {
		l = 150
	}
	form.Aboutme = strings.TrimSpace(form.Aboutme[0:l])
	l = len(form.Description)
	if l > 500 {
		l = 500
	}
	form.Description = strings.TrimSpace(form.Description[0:l])
	l = len(form.Fullname)
	if l > 150 {
		l = 150
	}
	form.Fullname = strings.TrimSpace(form.Fullname[0:l])

	form.Username = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(form.Username, "@")))

	links := &model.CreatorSocialMedia{
		Facebook:     form.Facebook,
		Twitter:      form.Twitter,
		PersonalSite: form.PersonalSite,
		Instagram:    form.Instagram,
		Whatsapp:     form.Whatsapp,
		Youtube:      form.Youtube,
	}
	validateSocialMedia(links)

	log.Printf("updating to %v", form)

	if !model.CheckUsername(form.Username) {
		apiError(ctx, fmt.Sprintf("%s is not a valid username", form.Username))
		return
	}

	tmp := creator.Username
	if creator, err = model.RetrieveCreatorByUsername(form.Username); err == nil && creator.Username != tmp {
		apiError(ctx, fmt.Sprintf("The username %s is already taken", form.Username))
		return
	} else {
		creator = ctx.Keys["creator"].(model.User)
	}

	if err = creator.UpdateBasicPageDetails(form.Fullname, form.Username, form.Aboutme, form.Description, links); err != nil {
		apiError(ctx, err)
		return
	}

	creator, err = model.RetrieveCreatorByID(creator.ID)
	util.ScrubbedPublicAPIJSON(ctx, creator, true)
}

//SetAvatar  updates the current user and sets the picture as their profile pic
func SetAvatar(ctx *gin.Context) {
	var err error
	var file model.CreatorFile
	creator := ctx.Keys["creator"].(model.User)

	id, ok := ctx.Params.Get("id")
	if !ok {
		apiError(ctx, fmt.Sprintf("id is not set"))
		return
	}

	var contentID primitive.ObjectID
	contentID, err = primitive.ObjectIDFromHex(id)

	if err != nil {
		apiError(ctx, fmt.Sprintf("invalid id"))
		return
	}

	//get file

	file, err = model.GetFileByID(contentID, nil, true)
	if err != nil {
		apiError(ctx, fmt.Sprintf("file does not exist %s", err))
		return
	}

	if file.Owner != creator.ID {
		apiError(ctx, fmt.Sprintf("you do not own this file"))
		return
	}

	if file.Type != "image" {
		apiError(ctx, fmt.Sprintf("file is not an image"))
		return
	}

	//lets do it :)
	err = creator.SetAvatar(file)
	if err != nil {
		apiError(ctx, fmt.Sprintf("failed to update your avatar"))
		return
	}
	creator.Profile.ProfilePhoto = file.ID
	util.ScrubbedPublicAPIJSON(ctx, file, true)
	return
}

//PublishPage publishes or unpublishes a creatorpage
func PublishPage(ctx *gin.Context) {

}

//AccountSummary generates the account summary
func AccountSummary(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"creator": creator,
	}, true)

}

//UpdatePageSocialMedia validates and updates a user's social media details
func UpdatePageSocialMedia(ctx *gin.Context) {

}

//UpdatePageLayout updates a creator's page layout
func UpdatePageLayout(ctx *gin.Context) {

}

type userBasicsForm struct {
	Fullname string `json:"fullname" form:"fullname" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Phone    string `json:"phone_number" form:"phone_number" binding:"required"`
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

//UpdateUserBasics updates the fullname, email and phone number
func UpdateUserBasics(ctx *gin.Context) {

	creator := ctx.Keys["creator"].(model.User)

	var form userBasicsForm
	var err error
	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, err)
		return
	}

	form.Fullname = strings.TrimSpace(form.Fullname)
	l := len(form.Fullname)
	if l > 150 {
		form.Fullname = form.Fullname[0:500]
	}

	form.Email = strings.ToLower(strings.TrimSpace(form.Email))
	form.Phone = strings.TrimSpace(form.Phone)

	if !isEmailValid(form.Email) {
		err = fmt.Errorf("%s is not a valid email address", form.Email)
		apiError(ctx, err)
		return
	}

	num, err := libphonenumber.Parse(form.Phone, "ZW")
	if err != nil {
		err = fmt.Errorf("Phone number is not valid. Only Zimbabwean phone numbers allowed")
		apiError(ctx, err)
		return
	}
	form.Phone = fmt.Sprintf("+%d%d", num.GetCountryCode(), num.GetNationalNumber())

	if err = creator.UpdateUserBasics(form.Fullname, form.Email, form.Phone); err != nil {
		apiError(ctx, err)
		return
	}

	go func() {
		if form.Email != creator.Email {
			//also notify of email change
			creator.Email = form.Email
			//resend verification code
			_, err := cron.SendEmailVerification(&creator, model.AccountSignup)
			if err != nil {
				log.Printf("Failed to send email verification for %v . Reason: %v", creator, err)
				//just continue, but something must be done. Ideally update and refect we failed to send the email
			}
		}

		if form.Phone != creator.PhoneNumber {
			creator.PhoneNumber = form.Phone
			_, err := cron.SendSMSVerification(&creator, model.ChangePhone)
			if err != nil {
				log.Printf("Failed to send sms verification for %v . Reason: %v", creator, err)
				util.ScrubbedPublicAPIJSON(ctx, gin.H{
					"error": fmt.Sprintf("Failed to send SMS but  to %s", creator.PhoneNumber),
				}, true)
				return
			}
		}
	}()

	util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "saved"}, true)
}

type pageConfigurablesForm struct {
	Item                 string  `json:"item" form:"item" `
	Price                float64 `json:"price" form:"price" `
	SupporterName        string  `json:"supporter" form:"supporter"`
	Thankyou             string  `json:"thanks" form:"thanks"`
	Googleanalytics      string  `json:"gacode" form:"gacode"`
	Supporters           bool    `json:"supporters" form:"supporters"`
	Subscribers          bool    `json:"subsactive" form:"subsactive"`
	SubscriptionMonths   int     `json:"subperiod" form:"subperiod"`
	SubscriptionPrice    float64 `json:"subprice" form:"subprice"`
	SubscriptionThankYou string  `json:"subthanks" form:"subthanks"`
	SubscriptionHeadline string  `json:"subheadline" form:"subheadline"`
}

var (
	allowedDonationItems  = map[string]int{"coffee": 0, "beer": 0, "lunch": 0, "pizza": 0, "chocolate": 0, "candy": 0, "ice cream": 0, "airtime": 0}
	allowedSupporterNames = map[string]int{"supporter": 0, "fan": 0, "client": 0, "member": 0, "backer": 0}
)

const defaultThankYouMessage = `
Thank you for supporting my hustle! 

Much love 
`

//UpdatePageConfigurables updates donation_item, price of donation item, supporters_name and thankyoumessage
func UpdatePageConfigurables(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	var form pageConfigurablesForm
	var err error
	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, err)
		return
	}

	l := len(form.Thankyou)
	if l > 500 {
		form.Thankyou = form.Thankyou[0:500]
	}

	price := form.Price

	if price < 0.5 {
		price = 0.5
	}

	if form.SupporterName == "" {
		form.SupporterName = "fan"
	}
	if _, ok := allowedSupporterNames[form.SupporterName]; !ok {
		apiError(ctx, fmt.Sprintf("supporter name is not supported"))
		return
	}

	if form.Item == "" {
		form.Item = "coffee"
	}
	if _, ok := allowedDonationItems[form.Item]; !ok {
		apiError(ctx, fmt.Sprintf("Donation item is not supported"))
		return
	}

	if form.Thankyou == "" {
		form.Thankyou = defaultThankYouMessage
	}

	page := creator.Page
	subs := creator.Subscriptions

	page.AllowSupporters = form.Supporters
	page.SupportersName = form.SupporterName
	page.DonationItemName = form.Item
	page.DonationItemUnitPrice = price
	page.ThankYouMessage = form.Thankyou

	if subs.Active != form.Subscribers {
		go cron.ScheduleNotifySubscriptionSettingsChange(creator, form.Subscribers)
	}
	subs.Active = form.Subscribers
	subs.Price = form.SubscriptionPrice
	subs.PeriodUnit = model.PeriodMonth
	subs.Period = int(math.Min(math.Max(float64(form.SubscriptionMonths), 12), 0))

	var subHeadline primitive.ObjectID
	if subHeadline, err = primitive.ObjectIDFromHex(form.SubscriptionHeadline); err != nil {
		subHeadline = creator.Profile.ProfilePhoto
	}
	subs.HeadlineImage = subHeadline
	subs.MaxSlots = 0
	subs.ThankYouMessage = form.SubscriptionThankYou

	if err = creator.UpdatePageConfigurables(page, subs); err != nil {
		apiError(ctx, err)
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "saved"}, true)
}

type formNotificationSettings struct {
	Login           bool `form:"login" json:"login" groups:"private"`
	WalletCredit    bool `form:"wallet_credit" json:"wallet_credit" groups:"private"`
	NewSubscriber   bool `form:"subscriber" json:"subscriber" groups:"private"`
	ServiceSchedule bool `form:"service_schedule" json:"service_schedule" groups:"private"`
}

//UpdateUserNotifications updates notification settings
func UpdateUserNotifications(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	var form formNotificationSettings
	var err error
	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, err)
		return
	}

	notifications := model.CreatorNotificationSettings{
		Login:           form.Login,
		WalletCredit:    form.WalletCredit,
		NewSubscriber:   form.NewSubscriber,
		ServiceSchedule: form.ServiceSchedule,
	}

	if err = creator.UpdateNotifications(notifications); err != nil {
		apiError(ctx, err)
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "saved"}, true)
}
