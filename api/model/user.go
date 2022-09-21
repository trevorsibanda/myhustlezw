package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	myhustleScheme   = os.Getenv("MYHUSTLE_SCHEME")
	myhustleEndpoint = os.Getenv("MYHUSTLE_DOMAIN")
)

type UserType string

var (
	UserTypeFan     UserType = "fan"
	UserTypeCreator UserType = "creator"
)

//User models a myhustle user
type User struct {
	ID                         primitive.ObjectID          `bson:"_id"  json:"_id" groups:"public"`
	Type                       UserType                    `bson:"type" json:"type" groups:"private"`
	Fullname                   string                      `bson:"fullname"  json:"fullname" groups:"public"`
	PhoneNumber                string                      `bson:"phone_number"  json:"phone_number" groups:"private,authenticated"`
	PhoneVerified              bool                        `bson:"phoneVerified"  json:"phoneVerified" groups:"private"`
	Email                      string                      `bson:"email"  json:"email" groups:"private,authenticated"`
	EmailVerified              bool                        `bson:"emailVerified"  json:"emailVerified" groups:"private"`
	PreviouslyVerified         bool                        `bson:"verified_ever"  json:"-" groups:"private"`
	VerificationPending        bool                        `bson:"verify_pending"  json:"verify_pending" groups:"private"`
	Created                    time.Time                   `bson:"created"  json:"created" groups:"private"`
	LastLogin                  time.Time                   `bson:"last_login"  json:"last_login" groups:"private"`
	WebNotification            WebNotification             `bson:"web_notifications" json:"web_notifications" groups:"private"`
	Username                   string                      `bson:"username"  json:"username" groups:"public"`
	LoggedIn                   bool                        `bson:"-" json:"logged_in" groups:"public"`
	IdentityVerified           bool                        `bson:"verified"  json:"verified" groups:"public"`
	Profile                    CreatorProfile              `bson:"profile"  json:"profile" groups:"public"`
	Page                       CreatorPageConfig           `bson:"page"  json:"page" groups:"public"`
	Notifications              CreatorNotificationSettings `bson:"notifications" json:"notifications" groups:"private"`
	Subscriptions              CreatorSubscription         `bson:"subscriptions" json:"subscriptions" groups:"public"`
	PayoutDetails              CreatorPayoutDetails        `bson:"payout" json:"payout,omitempty" groups:"private"`
	AcceptingPayments          bool                        `bson:"payments_active"  json:"payments_active" groups:"public"`
	PaymentsPassTransactionFee bool                        `bson:"payments_pass_txn_fee" json:"payments_pass_txn_fee" groups:"private"`
}

//CreatorSocialMedia stores a creator's socialmedia links
type CreatorSocialMedia struct {
	PersonalSite string `bson:"url" form:"url" json:"url" groups:"public"`
	Facebook     string `bson:"facebook" form:"facebook" json:"facebook" groups:"public"`
	Twitter      string `bson:"twitter" form:"twitter" json:"twitter" groups:"public"`
	Instagram    string `bson:"instagram" form:"instagram" json:"instagram" groups:"public"`
	Whatsapp     string `bson:"whatsapp" form:"whatsapp" json:"whatsapp" groups:"public"`
	Youtube      string `bson:"youtube" form:"youtube" json:"youtube" groups:"public"`
}

//CreatorProfile models a user's creator profile.
type CreatorProfile struct {
	CoverImage      primitive.ObjectID `bson:"cover_image"  json:"cover_image" groups:"private"`
	ProfilePhoto    primitive.ObjectID `bson:"profile_photo"  json:"profile_photo" groups:"private"`
	ProfilePhotoURL string             `bson:"-"  json:"profile_url" groups:"public"`
	AboutMe         string             `bson:"about_me"  json:"about_me" groups:"public"`
	Description     string             `bson:"description"  json:"description" groups:"public"`
	SocialMedia     CreatorSocialMedia `bson:"personal_website"  json:"personal_website" groups:"public"`
}

func (cp *CreatorProfile) SetProfilePicURL(ownerId primitive.ObjectID) {
	if ownerId.IsZero() {
		cp.ProfilePhotoURL = "/assets/img/user.png"
		return
	}
	cp.ProfilePhotoURL = CreatorFile{
		ID:      cp.ProfilePhoto,
		Owner:   ownerId,
		Type:    "image",
		Storage: string(LocalFile),
	}.PreviewURL(nil, true, 120, 120)
}

//CreatorPageConfig models a creator's page settings
type CreatorPageConfig struct {
	Published             bool               `bson:"published"  json:"published" groups:"public"`
	SupportersName        string             `bson:"supporter"  json:"supporter" groups:"public"`
	DonationItemName      string             `bson:"donation_item"  json:"donation_item" groups:"public"`
	DonationItemUnitPrice float64            `bson:"donation_item_unit_price"  json:"donation_item_unit_price" groups:"public"`
	ShowSocialMedia       bool               `bson:"socialmedia_info"  json:"socialmedia_info" groups:"public"`
	AllowSupporters       bool               `bson:"allow_supporters"  json:"allow_supporters" groups:"public"`
	ShowSupporterCounter  bool               `bson:"supporter_counter"  json:"supporter_counter" groups:"public"`
	ShowRecentSupporters  bool               `bson:"show_supporters"  json:"show_supporters" groups:"public"`
	FeaturedContent       primitive.ObjectID `bson:"featured_content"  json:"featured_content" groups:"public"`
	FeaturedURL           string             `bson:"-" json:"featured_url,omitempty" groups:"public"`
	GoogleAnalyticsCode   string             `bson:"google_analytics_code" json:"google_analytics_code" groups:"private"`
	ThankYouMessage       string             `bson:"thank_you_message" json:"thank_you_message" groups:"private"`
}

func (cpc *CreatorPageConfig) SetFeaturedURL(ownerId primitive.ObjectID) {
	if ownerId.IsZero() {
		cpc.FeaturedURL = "/assets/img/placeholder.png"
		return
	}
	cpc.FeaturedURL = CreatorFile{
		ID:      cpc.FeaturedContent,
		Owner:   ownerId,
		Type:    "image",
		Storage: string(LocalFile),
	}.PreviewURL(nil, true, 480, 480)
}

//CreatorSubscription stores a creator's subscription config
type CreatorSubscription struct {
	Active           bool               `bson:"active"  json:"active" groups:"public"`
	Price            float64            `bson:"price"  json:"price" groups:"public"`
	Period           int                `bson:"period"  json:"period" groups:"private"`
	PeriodUnit       SubscriptionPeriod `bson:"unit_period"  json:"unit_period" groups:"private"`
	MaxSlots         int                `bson:"max_slots"  json:"max_slots" groups:"private"`
	HeadlineImage    primitive.ObjectID `bson:"headline_image" json:"headline_image" groups:"public"`
	HeadlineImageURL string             `bson:"-" json:"url" groups:"public"`
	LastCount        int                `bson:"count" json:"count" groups:"public"`
	ThankYouMessage  string             `bson:"thank_you" json:"thank_you" groups:"private"`
}

//CreatorNotificationSettings stores a user's notification preferences
type CreatorNotificationSettings struct {
	Login           bool `bson:"login" json:"login" groups:"private"`
	WalletCredit    bool `bson:"wallet_credit" json:"wallet_credit" groups:"private"`
	NewSubscriber   bool `bson:"subscriber" json:"subscriber" groups:"private"`
	ServiceSchedule bool `bson:"service_schedule" json:"service_schedule" groups:"private"`
}

//URL creates your creator page url
func (creator User) URL() (url string) {
	return fmt.Sprintf("%s://%s/@%s", myhustleScheme, myhustleEndpoint, creator.Username)
}

//CheckUsername checks if a username is valid
func CheckUsername(username string) bool {
	pattern := "^([a-zA-Z0-9_]+)$"
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("failed to compile regexp expression in CheckUsername")
	}
	return re.Match([]byte(username))
}

//ProfilePhotoURL returns the url to the profile photo
func (creator User) ProfilePhotoURL() (url string) {
	url = CreatorFile{
		ID:      creator.Profile.ProfilePhoto,
		Owner:   creator.ID,
		Type:    "image",
		Storage: string(LocalFile),
	}.PreviewURL(nil, true, 240, 240)
	return
}

//CoverImageURL returns the url to the cover image
func (creator User) CoverImageURL() (url string) {
	url = CreatorFile{
		ID:      creator.Profile.CoverImage,
		Owner:   creator.ID,
		Type:    "image",
		Storage: string(LocalFile),
	}.PreviewURL(nil, true, 0, 0)
	return
}

//UpdateUserBasics updates the fullname, email and phone
func (creator User) UpdateUserBasics(fullname, email, phone string) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	phoneVerified := creator.PhoneVerified
	emailVerified := creator.EmailVerified

	if phone != creator.PhoneNumber {
		phoneVerified = false
	}

	if email != creator.Email {
		emailVerified = false
	}

	update := bson.M{
		"$set": bson.M{
			"email":         email,
			"emailVerified": emailVerified,
			"phoneVerified": phoneVerified,
			"fullname":      fullname,
			"phone_number":  phone,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//AnonymousUser creates an anonymous user
func AnonymousUser(email, phone, fullname string) (user User) {
	if fullname == "Someone" || fullname == "" {
		fullname = "Anonymous"
	}

	user = User{
		ID:            primitive.NilObjectID,
		Username:      "anonymous",
		Email:         email,
		EmailVerified: true,
		PhoneNumber:   phone,
		PhoneVerified: true,
		Fullname:      fullname,
		LoggedIn:      false,
	}
	user.Page.SetFeaturedURL(user.ID)
	return
}

//UpdatePageConfigurables updates donation_item, unit_price, thankyoumsg, price
func (creator User) UpdatePageConfigurables(page CreatorPageConfig, subs CreatorSubscription) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"page":          page,
			"subscriptions": subs,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//UpdatePageConfigurables updates donation_item, unit_price, thankyoumsg, price
func (creator User) UpdateNotifications(notifications CreatorNotificationSettings) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"notifications": notifications,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//UpdateBasicPageDetails does what it describes
func (creator User) UpdateBasicPageDetails(fullname, username, aboutMe, description string, links *CreatorSocialMedia) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	pf := creator.Profile
	pf.AboutMe = aboutMe
	pf.Description = description
	if links != nil {
		pf.SocialMedia = *links
	}

	update := bson.M{
		"$set": bson.M{
			"username": username,
			"fullname": fullname,
			"profile":  pf,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//RetrieveCreator retrieves a creator account for the specified  email, username and phone_number are used as keys
func RetrieveCreator(email, phoneNumber, username string) (creator User, err error) {
	keys := []bson.M{}
	if len(email) != 0 {
		keys = append(keys, bson.M{"email": strings.ToLower(email)})
	}
	if len(phoneNumber) != 0 {
		keys = append(keys, bson.M{"phone_number": phoneNumber})
	}
	if len(username) != 0 {
		keys = append(keys, bson.M{"username": strings.ToLower(username)})
	}
	if len(keys) == 0 {
		err = fmt.Errorf("Cannot retrieve user with empty email or phone number")
		return
	}

	err = creatorsCollection().FindOne(context.TODO(), bson.M{
		"$or": keys,
	}).Decode(&creator)
	creator.Page.SetFeaturedURL(creator.ID)
	creator.Profile.SetProfilePicURL(creator.ID)
	return
}

func (creator User) SetAvatar(image CreatorFile) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	pf := creator.Profile
	pf.ProfilePhoto = image.ID

	update := bson.M{
		"$set": bson.M{
			"profile": pf,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

func (creator User) SetCoverPicture(image CreatorFile) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	pf := creator.Profile
	pf.CoverImage = image.ID

	update := bson.M{
		"$set": bson.M{
			"profile": pf,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

func (creator User) SetFeaturedContent(image CreatorFile) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	pg := creator.Page
	pg.FeaturedContent = image.ID

	update := bson.M{
		"$set": bson.M{
			"page": pg,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

func UpdatePhoneNumberAndPaymentsToVerified(user User) (err error) {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"payments_active": true,
			"phoneVerified":   true,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

func UpdateEmailToVerified(user User) (err error) {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"emailVerified": true,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

func UpdatePhoneToVerified(user User) (err error) {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"phoneVerified": true,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

func (creator User) SetHeadlineSubscriptionImage(image CreatorFile) (err error) {
	filter := bson.M{
		"_id": creator.ID,
	}

	sb := creator.Subscriptions
	sb.HeadlineImage = image.ID

	update := bson.M{
		"$set": bson.M{
			"subscriptions": sb,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//UpdatePhoneNumber updates the creators phone number
func (user User) UpdatePhoneNumber(newNumber string) (err error) {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"phone_number":  newNumber,
			"phoneVerified": false,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//VerifyIdentity verifies the creators identity
func (user User) VerifyIdentity() (err error) {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"verified":       true,
			"verify_pending": false,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//UpdatePassword updates a user's password
func (user User) UpdatePassword(oldPassword, newPassword string) (err error) {
	filter := bson.M{
		"owner_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"password":    newPassword,
			"oldPassword": oldPassword,
			"lastChanged": time.Now(),
		},
	}

	_, err = securityCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//UpdateEmail updates the creators email
func (user User) UpdateEmail(newEmail string) (err error) {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"email":         newEmail,
			"emailVerified": false,
		},
	}

	_, err = creatorsCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//RetrieveCreatorByID retrieves a creator account given the ID
func RetrieveCreatorByID(id primitive.ObjectID) (creator User, err error) {
	filter := bson.M{
		"_id": id,
	}

	err = creatorsCollection().FindOne(context.TODO(), filter).Decode(&creator)
	creator.Page.SetFeaturedURL(creator.ID)
	creator.Profile.SetProfilePicURL(creator.ID)
	return
}

//RetrieveCreatorByUsername retrieves a creator given the username
func RetrieveCreatorByUsername(username string) (creator User, err error) {
	filter := bson.M{
		"username": username,
	}

	err = creatorsCollection().FindOne(context.TODO(), filter).Decode(&creator)
	creator.Page.SetFeaturedURL(creator.ID)
	creator.Profile.SetProfilePicURL(creator.ID)
	return
}

//RetrieveUser retrieves a fan/creator account for the specified user, using the email and phone_number or username
//all users have fan accounts. if none is set, then the user cannot exist.
func RetrieveUser(email, phoneNumber, username string) (creator User, err error) {
	creator, err = RetrieveCreator(email, phoneNumber, username)
	return
}

//NewCreator creates a new creator account for a specified user
func NewCreator(user User, username string, password string) (creator *User, sec *UserCredentials, err error) {
	user.ID = primitive.NewObjectID()
	user.Type = UserTypeCreator
	user.Username = username
	user.PreviouslyVerified = false
	user.Created = time.Now()
	user.AcceptingPayments = false
	user.LastLogin = time.Now()

	creator = &user

	_, err = creatorsCollection().InsertOne(context.TODO(), creator)
	if err != nil {
		log.Printf("Failed to create new creator account %v, Reason: %v", creator, err)
		return
	}

	sec, err = newUserCredentials(&user, password)
	if err != nil {
		log.Printf("Created creator account but failed to create security details, user: %v, Reason: %v", creator, err)
		filter := bson.D{{"_id", user.ID}}
		creatorsCollection().FindOneAndDelete(context.TODO(), filter)
	}
	creator.Page.SetFeaturedURL(creator.ID)
	creator.Profile.SetProfilePicURL(creator.ID)
	return
}

func DefaultCreatorProfile(skill string) (profile CreatorProfile) {
	profile = CreatorProfile{
		AboutMe:     fmt.Sprintf("is a %s", strings.ToLower(skill)),
		Description: "a content creator",
	}
	return
}

func DefaultCreatorPage() (page CreatorPageConfig) {
	page = CreatorPageConfig{
		SupportersName:  "supporter",
		AllowSupporters: true,
		Published:       true,
		ShowSocialMedia: true,
	}
	return
}

func DefaultMembershipsConfig() (m CreatorSubscription) {
	m = CreatorSubscription{
		Active:     true,
		Period:     12,
		PeriodUnit: PeriodMonth,
		Price:      1,
	}
	return
}
