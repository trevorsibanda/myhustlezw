package model

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/gosimple/slug"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	SingleVideoCampaign     string = "video"
	SingleAudioCampaign     string = "audio"
	SingleDownloadCampaign  string = "other"
	PhotobookCampaign       string = "image"
	MediaAlbumCampaign      string = "album"
	ServiceCampaign         string = "service"
	EmbeddedContentCampaign string = "embed"
)

var badCreatorURIS = map[string]interface{}{"buymeacoffee": nil, "thankyou": nil, "aboutme": nil, "subscribe": nil, "content": nil}

//Campaign models a creator campaign
type Campaign struct {
	ID              primitive.ObjectID `bson:"_id"  json:"_id" groups:"public"`
	URI             string             `bson:"uri"  json:"uri" groups:"public"`
	Owner           primitive.ObjectID `bson:"owner_id"  json:"owner_id" groups:"public"`
	Active          bool               `bson:"active"  json:"active" groups:"public"`
	Subscription    string             `bson:"subscription"  json:"subscription" groups:"public"`
	Title           string             `bson:"title"  json:"title" groups:"public"`
	Description     string             `bson:"description"  json:"description" groups:"public"`
	Created         time.Time          `bson:"created_at"  json:"created_at" groups:"public"`
	ExpiresAt       time.Time          `bson:"expires_at"  json:"expires_at" groups:"public"`
	LastUpdated     time.Time          `bson:"last_updated"  json:"last_updated" groups:"public"`
	Type            string             `bson:"type"  json:"type" groups:"public"`
	Price           float64            `bson:"price"  json:"price" groups:"public"` //price in cents
	Form            *ServiceForm       `bson:"service,omitempty"  json:"service,omitempty" groups:"public"`
	Log             string             `bson:"log"  json:"log" groups:"protected"`
	Preview         primitive.ObjectID `bson:"preview"  json:"preview" groups:"private"`
	PreviewImageURL string             `bson:"-"  json:"preview_url" groups:"public"`
	CanView         bool               `bson:"-" json:"can_view" groups:"public"`
	RemoteID        string             `bson:"remote_id,omitempty" json:"remote_id,omitempty" groups:"public"`
}

//ServiceForm stores information for service campaigns
type ServiceForm struct {
	ShowQuestion      bool     `bson:"show_question"  json:"show_question" groups:"public"`
	Question          string   `bson:"question"  json:"question" groups:"public"`
	ShowInstructions  bool     `bson:"show_instructions"  json:"show_instructions" groups:"public"`
	Instructions      string   `bson:"instructions"  json:"instructions" groups:"public"`
	ShowOptions       bool     `bson:"show_options"  json:"show_options" groups:"public"`
	OptionQuestion    string   `bson:"options_title"  json:"options_title" groups:"public"`
	OptionValues      []string `bson:"option_values"  json:"option_values" groups:"public"`
	ThankYouMessage   string   `bson:"thanks_message"  json:"thanks_message" groups:"private"`
	AllowExtraInfo    bool     `bson:"allow_extra_info"  json:"allow_extra_info" groups:"public"`
	ExtraInfoQuestion string   `bson:"extra_info_question"  json:"extra_info_question" groups:"public"`
	QuantityAvailable int      `bson:"quantity_available"  json:"quantity_available" groups:"private"`
}

//Fill fills a service form
func (form ServiceForm) Fill(email, phone, fullname, questionAnswer, selectedOption, extraInfoAnswer string, fulfilled bool) (filledForm FilledServiceForm) {
	filledForm = FilledServiceForm{
		Created:           time.Now(),
		Fulfilled:         fulfilled,
		Refunded:          false,
		Email:             email,
		Phone:             phone,
		Fullname:          fullname,
		Question:          form.Question,
		Answer:            questionAnswer,
		Instructions:      form.Instructions,
		OptionQuestion:    form.OptionQuestion,
		OptionValues:      form.OptionValues,
		SelectedOption:    selectedOption,
		ThankYouMessage:   form.ThankYouMessage,
		ExtraInfoQuestion: form.ExtraInfoQuestion,
		QuantityLeft:      form.QuantityAvailable,
	}
	return
}

func (service *Campaign) ServiceDecrementQuantity() (err error) {
	if service.Type != ServiceCampaign {
		err = fmt.Errorf("%s is not a service", service.ID.Hex())
		return
	}
	//update service to reduce items left by one
	service.Form.QuantityAvailable -= 1
	if service.Form.QuantityAvailable < 0 {
		service.Form.QuantityAvailable = 0
	}

	filter := bson.M{
		"_id": service.ID,
	}

	_, err = campaignsCollection().ReplaceOne(context.TODO(), filter, service)
	return
}

//Files returns all files for a user
func (campaign Campaign) Files(sessionKey *string, hasActiveSubscription bool) (files []CreatorFile, err error) {
	filter := bson.M{
		"campaign_id": campaign.ID,
	}

	var limit int64 = 1
	if hasActiveSubscription {
		limit = 100
	}

	options := options.FindOptions{
		Sort: bson.M{
			"created_at": -1,
		},
		Limit: &limit,
	}

	var cursor *mongo.Cursor

	cursor, err = creatorFilesCollection().Find(context.TODO(), filter, &options)
	if err == nil {
		err = cursor.All(context.TODO(), &files)
		for idx, file := range files {
			if hasActiveSubscription {
				file.URL = file.PreviewURL(sessionKey, hasActiveSubscription, 0, 0)
			} else {
				file.URL = file.PreviewURL(sessionKey, hasActiveSubscription, 480, 480)
			}
			if file.Type == "image" || file.Type == "video" {
				thumb := file.PreviewURL(sessionKey, hasActiveSubscription, 120, 120)
				file.Thumbnail = &thumb
			}
			files[idx] = file
		}
	}
	return
}

//Delete deletes the campaign
func (campaign Campaign) Delete() (err error) {
	campaign.Active = false
	campaign.Log = fmt.Sprintf("%s\n\nDeleted at %s", campaign.Log, time.Now().Format(time.UnixDate))
	filter := bson.M{
		"_id": campaign.ID,
	}

	_, err = campaignsCollection().ReplaceOne(context.TODO(), filter, campaign)
	return
}

//update the campaign
func (campaign Campaign) Update() (err error) {
	campaign.LastUpdated = time.Now()
	campaign.Log = fmt.Sprintf("%s\n\nUpdated at %s", campaign.Log, time.Now().Format(time.UnixDate))
	filter := bson.M{
		"_id": campaign.ID,
	}

	_, err = campaignsCollection().ReplaceOne(context.TODO(), filter, campaign)
	return
}

func (campaign Campaign) URL(creator string, supportID, unlockCode *primitive.ObjectID) string {
	if supportID == nil || unlockCode == nil {
		return fmt.Sprintf("%s://%s/%s/%s", myhustleScheme, myhustleEndpoint, creator, campaign.URI)
	}
	return fmt.Sprintf("%s://%s/v/%s/%s/%s/%s", myhustleScheme, myhustleEndpoint, creator, campaign.ID.Hex(), supportID.Hex(), unlockCode.Hex())
}

//PreviewURL generates an image preview url
func (campaign Campaign) PreviewURL(sessionKey *string, hasAccess bool, width, height int) (url string) {
	var file CreatorFile
	switch campaign.Type {
	case EmbeddedContentCampaign:

		file = CreatorFile{
			Type:     string(YoutubeEmbed),
			Storage:  string(YoutubeEmbed),
			Filename: campaign.RemoteID,
		}
	case SingleDownloadCampaign:
		file = CreatorFile{
			Type:     string(AssetFile),
			Storage:  string(AssetFile),
			Filename: "img/file.png",
		}
	default:
		file = CreatorFile{ID: campaign.Preview, Owner: campaign.Owner, Storage: string(LocalFile)}

	}
	url = previewURLGenerator(sessionKey, file, hasAccess, width, height)
	return
}

//CheckCampaignURI checks if the uri contains any illegal characters, returns true if valid
func CheckCampaignURI(uri string) bool {
	pattern := "^([a-zA-Z0-9_-]+)$"
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("failed to compile regexp expression in CreateCampaignURI")
	}
	return re.Match([]byte(uri))
}

//ComposeCampaignURI generates a uri given a campaign
func ComposeCampaignURI(c Campaign) string {
	check := slug.Make(c.Title)
	if _, ok := badCreatorURIS[check]; ok {
		check = fmt.Sprintf("%s%d", check, rand.Int())
	}
	for i := 0; i < 10; i++ {

		filter := bson.M{
			"$and": []bson.M{
				{"uri": check},
				{"owner_id": c.Owner},
			},
		}

		var campaign Campaign
		err := campaignsCollection().FindOne(context.TODO(), filter).Decode(&campaign)
		if err == nil && !campaign.ID.IsZero() {
			check = slug.Make(fmt.Sprintf("%s %d", c.Title, i))
		} else {
			return check
		}
	}

	return slug.Make(fmt.Sprintf("%s %d", c.Title, time.Now().Unix()))

}

//RetrieveCampaignByURI retrieves a creator campaign given the uri
func (user User) RetrieveCampaignByURI(sessionKey *string, uri string) (campaign Campaign, err error) {
	filter := bson.M{
		"$and": []bson.M{
			{"owner_id": user.ID},
			{"uri": uri},
			{"active": bson.M{"$ne": false}},
		},
	}

	err = campaignsCollection().FindOne(context.TODO(), filter).Decode(&campaign)
	if err == nil {
		if campaign.Subscription == "public" {
			campaign.PreviewImageURL = campaign.PreviewURL(nil, true, 480, 480)
		} else {
			campaign.PreviewImageURL = campaign.PreviewURL(sessionKey, false, 480, 480)
		}
	}
	return
}

//ListCampaigns retrieves a list of creator campaigns
func (user User) ListCampaigns() (campaigns []Campaign, err error) {
	campaigns = make([]Campaign, 0)
	filter := bson.M{
		"owner_id": user.ID,
		"active":   bson.M{"$ne": false},
	}

	opts := options.FindOptions{}

	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := campaignsCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	err = cursor.All(context.TODO(), &campaigns)
	return
}

//ListRecentCampaigns list all recent campaigns for a creator
func (user User) ListRecentCampaigns(sessionKey *string, maxItems int64, skip int64, hasActiveSubscriptions bool) (list []Campaign, err error) {
	list = make([]Campaign, 0)
	var cursor *mongo.Cursor
	filter := bson.M{
		"owner_id": user.ID,
		"active":   bson.M{"$ne": false},
	}
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &maxItems,
	}

	opts.SetSort(bson.M{"created_at": -1})

	cursor, err = campaignsCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	cursor.All(context.TODO(), &list)
	for idx, campaign := range list {
		if campaign.Subscription == "fans" || campaign.Subscription == "pay_per_view" {
			campaign.PreviewImageURL = campaign.PreviewURL(sessionKey, hasActiveSubscriptions, 480, 480)
		} else {
			campaign.PreviewImageURL = campaign.PreviewURL(sessionKey, true, 480, 480)
		}
		campaign.CanView = hasActiveSubscriptions
		list[idx] = campaign
	}

	return
}

//GetCampaignsFiles returns a list of campaign files given a list of campaigns
func GetCampaignsFiles(campaigns []Campaign) (files map[primitive.ObjectID][]CreatorFile, err error) {
	files = make(map[primitive.ObjectID][]CreatorFile)
	var cursor *mongo.Cursor

	var ids bson.A
	for _, campaign := range campaigns {
		files[campaign.ID] = make([]CreatorFile, 0)
		ids = append(ids, bson.M{"campaign_id": campaign.ID})
	}

	opts := options.FindOptions{}
	opts.SetSort(bson.M{"created_at": -1})
	filter := bson.M{
		"$or": ids,
	}

	cursor, err = creatorFilesCollection().Find(context.TODO(), filter, &opts)
	if err != nil {
		return
	}
	println(err)
	defer cursor.Close(context.TODO())

	list := make([]CreatorFile, 0)
	cursor.All(context.TODO(), &list)

	for _, file := range list {
		files[file.ID] = append(files[file.ID], file)
	}

	return
}

//NewServiceCampaign creates a service campaign with the given options
func (user User) NewServiceCampaign(title, description string, price float64, form ServiceForm, preview *CreatorFile) (c *Campaign, err error) {
	c = &Campaign{
		ID:           primitive.NewObjectID(),
		URI:          "generate-uri-here",
		Owner:        user.ID,
		Active:       true,
		Preview:      preview.ID,
		Subscription: "public",
		Type:         ServiceCampaign,
		Price:        price,
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		Title:        title,
		Description:  description,
		Form:         &form,
	}

	c.URI = ComposeCampaignURI(*c)

	_, err = campaignsCollection().InsertOne(context.TODO(), c)
	return
}

//NewVideoCampaign creates a new video campaign
func (user User) NewVideoCampaign(title, description, subscription string, price float64, video *CreatorFile) (c *Campaign, err error) {
	c = &Campaign{
		ID:           primitive.NewObjectID(),
		URI:          "",
		Owner:        user.ID,
		Active:       true,
		Subscription: subscription,
		Preview:      video.ID,
		Type:         SingleVideoCampaign,
		Price:        price,
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		Title:        title,
		Description:  description,
	}

	c.URI = ComposeCampaignURI(*c)

	_, err = campaignsCollection().InsertOne(context.TODO(), c)
	if err == nil {
		setCampaignForFiles(c, []*CreatorFile{video})
	}
	return
}

//NewAudioCampaign creates a new video campaign
func (user User) NewAudioCampaign(title, description, subscription string, price float64, audio *CreatorFile) (c *Campaign, err error) {
	c = &Campaign{
		ID:           primitive.NewObjectID(),
		URI:          "",
		Owner:        user.ID,
		Active:       true,
		Subscription: subscription,
		Preview:      audio.ID,
		Type:         SingleAudioCampaign,
		Price:        price,
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		Title:        title,
		Description:  description,
	}

	c.URI = ComposeCampaignURI(*c)

	_, err = campaignsCollection().InsertOne(context.TODO(), c)
	if err == nil {
		setCampaignForFiles(c, []*CreatorFile{audio})
	}
	return
}

//NewPhotoAlbum creates a new photobook campaign
func (user User) NewPhotobook(title, description, subscription string, price float64, images []*CreatorFile) (c *Campaign, e error) {
	c = &Campaign{
		ID:           primitive.NewObjectID(),
		URI:          "",
		Owner:        user.ID,
		Active:       true,
		Preview:      images[0].ID,
		Subscription: subscription,
		Type:         PhotobookCampaign,
		Price:        price,
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		Title:        title,
		Description:  description,
	}

	c.URI = ComposeCampaignURI(*c)

	_, err := campaignsCollection().InsertOne(context.TODO(), c)
	if err == nil {
		setCampaignForFiles(c, images)
	}
	return
}

//NewEmbedCampaign creates a new embedded content campaign
func (user User) NewEmbedCampaign(campaign Campaign) (c *Campaign, err error) {
	c = &Campaign{
		ID:           primitive.NewObjectID(),
		URI:          "",
		Owner:        campaign.Owner,
		Active:       campaign.Active,
		Subscription: campaign.Subscription,
		Type:         EmbeddedContentCampaign,
		Price:        campaign.Price,
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		ExpiresAt:    campaign.ExpiresAt,
		Title:        campaign.Title,
		Description:  campaign.Description,
		RemoteID:     campaign.RemoteID,
	}

	c.URI = ComposeCampaignURI(*c)
	_, err = campaignsCollection().InsertOne(context.TODO(), c)
	return
}

//NewSingleDownload creates a general download campaign
func (user User) NewSingleDownload(title, description, subscription string, price float64, file *CreatorFile) (c *Campaign, err error) {
	c = &Campaign{
		ID:           primitive.NewObjectID(),
		URI:          "",
		Owner:        user.ID,
		Active:       true,
		Subscription: subscription,
		Type:         SingleDownloadCampaign,
		Price:        price,
		Created:      time.Now(),
		LastUpdated:  time.Now(),
		Title:        title,
		Description:  description,
	}

	c.URI = ComposeCampaignURI(*c)

	_, err = campaignsCollection().InsertOne(context.TODO(), c)
	if err == nil {
		setCampaignForFiles(c, []*CreatorFile{file})
	}
	return
}

func setCampaignForFiles(campaign *Campaign, files []*CreatorFile) (err error) {
	for _, file := range files {
		err = file.SetCampaign(campaign)
		if err != nil {
			_, err = campaignsCollection().DeleteOne(context.TODO(), bson.M{"_id": campaign.ID})
			if err != nil {
				log.Println("failed to update campaign file to reflect new campaign and failed to delete fileless campaign.")
			} else {
				err = fmt.Errorf("created campaign but failed to update file")
			}
			return
		}
	}
	return
}

//CampaignByID retrieves a campaign given the id
func CampaignByID(id primitive.ObjectID) (campaign Campaign, err error) {
	filter := bson.M{
		"_id":    id,
		"active": bson.M{"$ne": false},
	}

	err = campaignsCollection().FindOne(context.TODO(), filter).Decode(&campaign)
	return
}
