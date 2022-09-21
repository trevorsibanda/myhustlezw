package api

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/storage"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ListCampaigns lists all the creators campaigns
func ListCampaigns(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	//todo, allow filtering

	campaigns, err := creator.ListCampaigns()
	if err != nil {
		//log the error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve campaigns",
		})
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, campaigns, true)
}

type createServiceForm struct {
	Title        string  `form:"title" binding:"required"`
	Type         string  `form:"type" binding:"required"`
	Description  string  `form:"description"`
	Instructions string  `form:"instructions" binding:"required"`
	Question     string  `form:"question" binding:"required"`
	Price        float64 `form:"price" binding:"required"`
	Quantity     string  `form:"quantity" `
	Thankyou     string  `form:"thankyou" binding:"required"`
	Preview      string  `form:"preview" binding:"required"`
	Expires      string  `form:"expires" `
	Subscription string  `form:"subscription" binding:"required"`
}

type createDigitalForm struct {
	Title        string   `form:"title" binding:"required"`
	Type         string   `form:"type" binding:"required"`
	Description  string   `form:"description" binding:"required"`
	Price        float64  `form:"price"`
	Content      []string `form:"content" binding:"required"`
	Expires      string   `form:"expires" binding:"required"`
	Subscription string   `form:"subscription" binding:"required"`
}

//For now all forms expire in 6 months from creation date
func calculateFormExpires(fromNow string) (when time.Time) {
	now := time.Now()
	when = now.AddDate(0, 6, 0)
	return
}

//CreateNewCampaign creates a new campaign
func CreateNewCampaign(ctx *gin.Context) {
	tpeS, _ := ctx.Params.Get("type")
	tpeS = strings.ToLower(tpeS)
	switch tpeS {
	case "service":
		CreateServiceCampaign(ctx)
	case "youtube_embed":
		CreateEmbedCampaign(ctx)
	case "embed":
		CreateEmbedCampaign(ctx)
	case "content":
		CreateDigitalCampaign(ctx)
	default:
		apiError(ctx, fmt.Sprintf("Cannot create %s campaign", tpeS))
	}
}

//CreateDigitalCampaign creates a new service campaign
func CreateDigitalCampaign(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)

	var form createDigitalForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid request form %v", err))
		return
	}

	files := []*model.CreatorFile{}

	for _, id := range form.Content {
		fileID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			apiError(ctx, fmt.Sprintf("preview image does not exist %v", err))
			return
		}
		content, err := model.GetFileByID(fileID, &(session.PrivateKey), true)
		if err != nil {
			apiError(ctx, fmt.Sprintf("provided preview image does not exist: %v", err))
			return
		}
		if content.Owner == creator.ID {
			files = append(files, &content)
		}
	}

	if len(files) == 0 {
		apiError(ctx, "At least one file is required to create a campaign")
		return
	}

	var tpe string = form.Type

	if tpe != model.SingleVideoCampaign && tpe != model.SingleAudioCampaign && tpe != model.PhotobookCampaign && tpe != model.SingleDownloadCampaign {
		apiError(ctx, fmt.Sprintf("Invalid upload type %s", tpe))
		return
	}
	price := form.Price

	switch form.Subscription {
	case "public":
		price = 0.00
	case "private":
		form.Subscription = "fans"
	case "fans":
	case "pay_per_view":
		if price < 0.50 {
			apiError(ctx, fmt.Sprintf("minimum view price is $1.00"))
			return
		}
		price = math.Min(price, 100)
	default:
		apiError(ctx, "Invalid subscription group")
		return
	}

	c := &model.Campaign{
		URI:          "generate-uri-here",
		Owner:        creator.ID,
		Subscription: form.Subscription,
		Title:        form.Title,
		Price:        price,
		Type:         tpe,
		Created:      time.Now(),
		ExpiresAt:    calculateFormExpires(form.Expires),
		Description:  form.Description,
	}

	//check preview image and only use it if its a service preview
	var err error
	var campaign *model.Campaign

	switch c.Type {
	case model.SingleVideoCampaign:
		campaign, err = creator.NewVideoCampaign(c.Title, c.Description, c.Subscription, c.Price, files[0])
	case model.SingleAudioCampaign:
		campaign, err = creator.NewAudioCampaign(c.Title, c.Description, c.Subscription, c.Price, files[0])
	case model.PhotobookCampaign:
		campaign, err = creator.NewPhotobook(c.Title, c.Description, c.Subscription, c.Price, files)
	case model.SingleDownloadCampaign:
		campaign, err = creator.NewSingleDownload(c.Title, c.Description, c.Subscription, c.Price, files[0])
	default:
		apiError(ctx, "Campaign type not yet supported")
		return
	}

	if err != nil {
		apiError(ctx, "Failed to save in database")
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, campaign, true)

}

//CreateEmbedCampaign creates an embeded content campaign
func CreateEmbedCampaign(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	var form createDigitalForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid request form:\n %v", err))
		return
	}

	var youtubeVideo string
	if len(form.Content) == 0 {
		apiError(ctx, "At least one ID is required to create a campaign")
		return
	}

	youtubeVideo = strings.TrimSpace(form.Content[0])

	//TODO: check if youtube code is valid

	tpe := model.EmbeddedContentCampaign
	price := 0.00

	switch form.Subscription {
	case "public":
		price = 0.00
	default:
		apiError(ctx, fmt.Sprintf("invalid subscription group"))
		return
	}

	c := &model.Campaign{
		URI:          "generate-uri-here",
		Owner:        creator.ID,
		Subscription: form.Subscription,
		Title:        form.Title,
		Price:        price,
		Active:       true,
		Type:         tpe,
		Created:      time.Now(),
		ExpiresAt:    calculateFormExpires(form.Expires),
		Description:  form.Description,
		RemoteID:     youtubeVideo,
	}

	//check preview image and only use it if its a service preview
	var err error
	var campaign *model.Campaign

	campaign, err = creator.NewEmbedCampaign(*c)

	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save in database"))
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, campaign, true)

}

//CreateServiceCampaign creates a new service campaign
func CreateServiceCampaign(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)

	var form createServiceForm

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid request form %v", err))
		return
	}

	var fileID primitive.ObjectID

	fileID, err := primitive.ObjectIDFromHex(form.Preview)
	if err != nil {
		apiError(ctx, fmt.Sprintf("preview image does not exist"))
		return
	}

	quantity, _ := strconv.Atoi(form.Quantity)

	serviceForm := model.ServiceForm{
		ShowQuestion:      true,
		Question:          form.Question,
		ShowInstructions:  true,
		Instructions:      form.Instructions,
		ShowOptions:       false,
		ThankYouMessage:   form.Thankyou,
		QuantityAvailable: int(quantity),
		AllowExtraInfo:    false,
	}

	price := form.Price
	if price < 0.50 {
		apiError(ctx, fmt.Sprintf("minimum service price is $0.50"))
		return
	}

	c := &model.Campaign{
		URI:          "generate-uri-here",
		Owner:        creator.ID,
		Subscription: "public",
		Title:        form.Title,
		Price:        price,
		Type:         "service",
		Created:      time.Now(),
		ExpiresAt:    calculateFormExpires(form.Expires),
		Description:  form.Description,
		Form:         &serviceForm,
	}

	preview, err := model.GetFileByID(fileID, &(session.PrivateKey), true)
	if err != nil {
		apiError(ctx, fmt.Sprintf("provided preview image does not exist"))
		return
	}
	if preview.Owner != creator.ID {
		apiError(ctx, fmt.Sprintf("you do not own this file"))
		return
	}
	if preview.Type != "video" && preview.Type != "image" {
		apiError(ctx, fmt.Sprintf("only videos and images can be used as previews. Found %s", preview.Type))
		return
	}

	campaign, err := creator.NewServiceCampaign(c.Title, c.Description, c.Price, serviceForm, &preview)
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save in database"))
		return
	}

	preview.Campaign = campaign.ID
	if err = preview.SaveChanges(); err != nil {
		log.Printf("Failed to update preview image %v %v %v", creator, campaign, preview)
	}

	util.ScrubbedPublicAPIJSON(ctx, campaign, true)

}

//GetDetailedCampaign returns a detailed campaign
func GetDetailedCampaign(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)

	var campaign model.Campaign
	var files []model.CreatorFile
	var orders []model.FilledServiceForm
	var supporters []model.CreatorSupport
	var stats gin.H

	var err error

	var id primitive.ObjectID
	if idString, ok := ctx.Params.Get("id"); !ok {
		apiError(ctx, fmt.Sprintf("No ID passed"))
		return
	} else if id, err = primitive.ObjectIDFromHex(idString); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid ID passed"))
		return
	}

	if campaign, err = model.CampaignByID(id); err != nil {
		apiError(ctx, err)
		return
	}

	if campaign.Owner != creator.ID {
		apiError(ctx, fmt.Sprintf("cannot access campaigns that do not belong to you"))
		return
	}

	if files, err = campaign.Files(&(session.PrivateKey), true); err != nil {
		apiError(ctx, fmt.Sprintf("failed to retrieve files for campaign %s", err))
		return
	}

	if len(files) == 0 && campaign.Type == "service" {
		file, _ := model.GetFileByID(campaign.Preview, &(session.PrivateKey), true)
		files = append(files, file)
	}

	for idx, file := range files {
		file.URL, _ = storage.GetImageURL(&session.PrivateKey, file, 0, 0, false)
		files[idx] = file
	}

	if supporters, err = campaign.ListRecentSupporters(20, 0); err != nil {
		apiError(ctx, fmt.Sprintf("failed to retrieve recent supporters %s", err))
		return
	}

	if campaign.Type != model.ServiceCampaign {
		orders = make([]model.FilledServiceForm, 0)
	} else if orders, err = campaign.ListRecentSubmittedForms(20, 0); err != nil {
		apiError(ctx, fmt.Sprintf("failed to retrieve recent orders%s", err))
		return
	}

	stats = gin.H{
		"views":     0,
		"downloads": 0,
		"orders":    0,
		"comments":  []string{},
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"campaign":   campaign,
		"files":      files,
		"supporters": supporters,
		"stats":      stats,
		"orders":     orders,
	}, true)

}

type updateCampaignForm struct {
	Title        string            `json:"title" binding:"required" form:"title" `
	Price        float64           `json:"price" binding:"required" form:"price"`
	Description  string            `json:"description" form:"description"`
	Subscription string            `json:"subscription" form:"subscription"`
	Preview      string            `json:"preview" form:"preview"`
	Service      model.ServiceForm `json:"service" form:"service"`
}

//UpdateCampaign updates a campaign
func UpdateCampaign(ctx *gin.Context) {

	creator := ctx.Keys["creator"].(model.User)

	var campaign model.Campaign

	var err error

	var id primitive.ObjectID
	if idString, ok := ctx.Params.Get("id"); !ok {
		apiError(ctx, "No ID passed")
		return
	} else if id, err = primitive.ObjectIDFromHex(idString); err != nil {
		apiError(ctx, "Invalid ID passed")
		return
	}

	var form updateCampaignForm
	if err = ctx.BindJSON(form); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to parse form with error %v", err))
		return
	}

	if campaign, err = model.CampaignByID(id); err != nil {
		apiError(ctx, err)
		return
	}

	if campaign.Owner != creator.ID {
		apiError(ctx, "Cannot access campaigns that do not belong to you")
		return
	}

	price := campaign.Price
	if form.Price < 0.5 || form.Price > 100 {
		apiError(ctx, "Price is outside acceptable bounds")
		return
	}

	sub := campaign.Subscription
	if campaign.Type != "service" && (form.Subscription != "public" && form.Subscription != "pay_per_view" && form.Subscription != "fans") {
		apiError(ctx, "Unkown subscription type")
		return
	}

	campaign.Title = strings.TrimSpace(form.Title)
	campaign.Description = strings.TrimSpace(form.Description)
	campaign.Subscription = sub
	campaign.Price = price

	if campaign.Type == "service" && (form.Service.Question == "" || campaign.Form == nil) {
		apiError(ctx, "Service has incomplete request fields")
		return
	}
	if campaign.Type == "service" {
		campaign.Form = &form.Service
	}

	//todo: allow updating preview image
	//files, _ = campaign.Files(&(session.PrivateKey), true)

	if err = campaign.Update(); err != nil {
		apiError(ctx, "Failed to save your changes")
	}

	util.ScrubbedPublicAPIJSON(ctx, campaign, true)

}

func DeleteCampaign(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)

	var campaign model.Campaign
	var files []model.CreatorFile

	var err error

	var id primitive.ObjectID
	if idString, ok := ctx.Params.Get("id"); !ok {
		apiError(ctx, fmt.Sprintf("No ID passed"))
		return
	} else if id, err = primitive.ObjectIDFromHex(idString); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid ID passed"))
		return
	}

	if campaign, err = model.CampaignByID(id); err != nil {
		apiError(ctx, err)
		return
	}

	if campaign.Owner != creator.ID {
		apiError(ctx, fmt.Sprintf("cannot access campaigns that do not belong to you"))
		return
	}

	//todo: for services first check that there are not pending orders to be fulfilled, then delete

	files, err = campaign.Files(&(session.PrivateKey), true)

	go func() {
		for _, file := range files {
			file.Delete()
		}
	}()

	campaign.Delete()

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "deleted",
	}, true)

}

//ListRecentCampaignSupporters list recent supporters up to :max_per_page
func ListRecentCampaignSupporters(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	var maxItems int64 = 10
	var skip int64 = 0

	maxPage, ok := ctx.Params.Get("max_per_page")
	if ok {
		tmp, err := strconv.Atoi(maxPage)
		if err != nil {
			maxItems = int64(tmp)
		}
	}
	maxPage, ok = ctx.Params.Get("skip")
	if ok {
		tmp, err := strconv.Atoi(maxPage)
		if err != nil {
			skip = int64(tmp)
		}
	}

	idString, _ := ctx.Params.Get("id")
	id, _ := primitive.ObjectIDFromHex(idString)
	if campaign, err := model.CampaignByID(id); err != nil {
		apiError(ctx, "Content does not exist")
		return
	} else if campaign.Owner != creator.ID {
		apiError(ctx, fmt.Sprintf("cannot access campaigns that do not belong to you"))
		return
	}

	//todo add suppprt for cursors

	supporters, err := creator.ListRecentSupporters(maxItems, skip, &id, true)
	if err != nil {
		apiError(ctx, fmt.Sprintf("failed to retrieve supporters %s", err))
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"skip":       skip,
		"page_size":  maxItems,
		"supporters": supporters,
	}, true)

}
