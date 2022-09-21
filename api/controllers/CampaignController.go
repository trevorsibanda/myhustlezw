package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/storage"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var defaultPreviewImage model.CreatorFile

//PublicDownloadCampaignFile downloads a file from a campaign
func PublicDownloadCampaignFile(ctx *gin.Context) {
	apiError(ctx, "Not implemented")
}

//PublicRequestService handles a request to a service
func PublicRequestService(ctx *gin.Context) {
	apiError(ctx, "Not implemented")
}

//PublicGetServiceRequest returns a service request
func PublicGetServiceRequest(ctx *gin.Context) {
	apiError(ctx, "Not implemented")
}

type serveWriter = func(ctx *gin.Context, data interface{}, loggedIn bool)

//CampaignPageView renders a campaign page based on the writer
func PublicGetCampaign(writer serveWriter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creator := ctx.Keys["creator"].(model.User)
		session := ctx.Keys["session"].(sessions.VisitorSession)
		uri := ctx.Params.ByName("id")

		//todo: check if campaign exists using id
		if !model.CheckCampaignURI(uri) {
			apiError(ctx, "Invalid content ID passed")
			return
		}

		campaign, err := creator.RetrieveCampaignByURI(&(session.PrivateKey), uri)
		if err != nil {
			apiError(ctx, "Content does not exist or has been deleted.")
			return
		}

		if !campaign.Active {
			apiError(ctx, "Content has been deleted.")
			return
		}

		canViewCreator, sub := canViewCreatorSubscriberContent(creator, session, &campaign, false)

		type handler = func(*gin.Context, *Page, bool, *model.CreatorSupport, *sessions.VisitorSession, *model.User, *model.Campaign) gin.H

		var fn handler
		page := newPage("")
		page.SetPageURL(campaign.URL(creator.Username, nil, nil))

		page.SetTwitterUsername(creator.Profile.SocialMedia.Twitter)
		page.SetPageAuthorURL(creator.URL())
		page.SetImage(campaign.PreviewImageURL)

		if campaign.Subscription == "public" {
			desc := campaign.Description
			if len(desc) > 150 {
				desc = campaign.Description[:150]
			}
			page.SetDescription(desc)
		} else if campaign.Subscription == "fans" {
			desc := fmt.Sprintf("Subscribe to @%s for only %s/%s for 3months to access this content exclusive to %ss only",
				creator.Username, cron.FormatAsUSD(campaign.Price), cron.FormatAsZWL(cron.USDToZWL(campaign.Price)), creator.Page.SupportersName)
			page.SetDescription(desc)
			page.SetImage(campaign.PreviewURL(nil, false, 480, 480))
		} else if campaign.Subscription == "pay_per_view" {
			desc := fmt.Sprintf("Pay %s/%s to access this private post by @%s",
				cron.FormatAsUSD(campaign.Price), cron.FormatAsZWL(cron.USDToZWL(campaign.Price)), creator.Username)
			page.SetDescription(desc)
			page.SetImage(campaign.PreviewURL(nil, false, 480, 480))
		}
		switch campaign.Type {
		case model.SingleVideoCampaign:
			page.SetImageAlt("video preview")
			fn = singleVideoCampaign
			page.SetTitle(fmt.Sprintf("[VIDEO] %s by @%s", campaign.Title, creator.Username))

		case model.SingleAudioCampaign:
			fn = singleVideoCampaign
			page.SetImageAlt("audio track or album cover")
			page.SetTitle(fmt.Sprintf("[AUDIO] %s by @%s", campaign.Title, creator.Username))
		case model.SingleDownloadCampaign:
			fn = downloadFileCampaign
			page.SetImageAlt("download file cover")
			page.SetTitle(fmt.Sprintf("[DOWNLOAD] %s by @%s", campaign.Title, creator.Username))
		case model.ServiceCampaign:
			fn = serviceCampaign
			page.SetImageAlt("service alt image")
			page.SetTitle(fmt.Sprintf("[SERVICE] %s offered by @%s", campaign.Title, creator.Username))
		case model.MediaAlbumCampaign:
			fn = albumCampaign
			page.SetImageAlt("photo album preview")
			page.SetTitle(fmt.Sprintf("[PHOTOS] %s by @%s", campaign.Title, creator.Username))
		case model.PhotobookCampaign:
			fn = photobookCampaign
			page.SetImageAlt("photos preview")
			page.SetTitle(fmt.Sprintf("[PHOTOS] %s by @%s", campaign.Title, creator.Username))
		case model.EmbeddedContentCampaign:
			fn = embedCampaign
			page.SetImageAlt("remote preview image")
			page.SetTitle(fmt.Sprintf("[VIDEO] %s by @%s", campaign.Title, creator.Username))
		default:
			apiError(ctx, fmt.Sprintf("Internal server error. Campaign has incorrect database type"))
			return
		}

		data := fn(ctx, &page, canViewCreator, &sub, &session, &creator, &campaign)
		data["page"] = &page
		data["user"] = session.User
		data["content"] = campaign
		data["creator"] = creator
		data["featured"], _ = model.GetFileByID(creator.Page.FeaturedContent, nil, true)
		data["subscription"] = sub
		data["content"] = campaign

		feed, _ := session.User.GenerateUserFeed(0, model.FeedViewDiscover, model.NoFilter, session.PaidContent, nil, nil, nil)

		subs, _ := session.User.GetSubscriptions("active")
		data["recommendations"] = []model.Campaign{}
		for _, item := range feed {
			canViewCreator := false
			for _, s := range subs {
				if s.Creator == item.Content.Owner {
					canViewCreator = true
					break
				}
			}
			canView, _ := CanViewCreatorSubscriberContent(item.Creator, session, &(item.Content), canViewCreator)
			item.Content.CanView = canView
			item.Content.PreviewImageURL = item.Content.PreviewURL(&session.PrivateKey, canView, 480, 480)
			data["recommendations"] = append(data["recommendations"].([]model.Campaign), item.Content)
		}

		data["allowed"] = canViewCreator

		b := session.User.LoggedIn

		writer(ctx, data, b)
	}
}

func dashboardFile(file string) string {
	return "./static/dashboard/build/" + file
}

//writer adds meta tags to the html page for all campaigns
//this is also cached
func GenerateCachedCampaignPage(ctx *gin.Context, data interface{}, loggedIn bool) {
	file, err := os.ReadFile(dashboardFile("index.html"))
	if err != nil {
		apiError(ctx, "Failed to render page.")
		return
	}
	page := data.(gin.H)["page"].(*Page)
	meta := util.GenerateMeta(page.Dict)
	editted := strings.Replace(string(file), "<pagemeta></pagemeta>", meta, 1)
	withPreGen := util.GeneratePageWithData("content", editted, data)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(withPreGen))
}

func singleVideoCampaign(ctx *gin.Context, page *Page, canViewCreator bool, sub *model.CreatorSupport, session *sessions.VisitorSession, creator *model.User, campaign *model.Campaign) (data gin.H) {
	//load preview media if available
	var video model.CreatorFile
	var recentSupporters []model.CreatorSupport

	//should just findone instead
	if files, err := campaign.Files(&(session.PrivateKey), canViewCreator); err == nil || len(files) != 0 {
		video = files[0]
	} else {
		//file was deleted
		apiError(ctx, fmt.Sprintf("Internal server error. Error retrieving files for video %s", campaign.ID))
	}
	streamURL, err := storage.GetVideoStreamURL(video, canViewCreator)
	if err != nil {
		apiError(ctx, fmt.Sprintf("Internal server error. Error retrieving stream URL for video %s", campaign.ID))
	}

	if session.User.ID == creator.ID {
		session.User.LoggedIn = true

		recentSupporters, err = creator.ListRecentSupporters(8, 0, &(campaign.ID), true)
		if err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		}
	}

	data = gin.H{
		"page":         page,
		"creator":      creator,
		"subscription": sub,
		"content":      campaign,
		"supporters":   recentSupporters,
		"video":        video,
		"stream_url":   streamURL,
	}

	return

}

func downloadFileCampaign(ctx *gin.Context, page *Page, canViewCreator bool, sub *model.CreatorSupport, session *sessions.VisitorSession, creator *model.User, campaign *model.Campaign) (data gin.H) {
	//load preview media if available
	var err error
	var files []model.CreatorFile
	var recentSupporters []model.CreatorSupport

	if files, err = campaign.Files(&(session.PrivateKey), canViewCreator); err != nil {
		//show deleted page
		apiError(ctx, fmt.Sprintf("Internal server error. Error retrieving files for campaign %s", campaign.Title))
		return
	}

	if !session.User.ID.IsZero() {
		session.User.LoggedIn = true

		recentSupporters, err = creator.ListRecentSupporters(8, 0, &(campaign.ID), true)
		if err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		}
	}

	data = gin.H{
		"supporters": recentSupporters,
		"file":       files[0],
	}
	return
}

func serviceCampaign(ctx *gin.Context, page *Page, canViewCreator bool, sub *model.CreatorSupport, session *sessions.VisitorSession, creator *model.User, campaign *model.Campaign) (data gin.H) {

	var recentSupporters []model.CreatorSupport

	//load preview media if available
	preview, err := model.GetFileByID(campaign.Preview, &(session.PrivateKey), true)
	if err != nil {
		preview = defaultPreviewImage
	}

	if session.User.ID == creator.ID {
		session.User.LoggedIn = true

		recentSupporters, err = creator.ListRecentSupporters(8, 0, &(campaign.ID), true)
		if err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		}
	}

	data = gin.H{
		"preview":    preview,
		"supporters": recentSupporters,
	}
	return
}

func albumCampaign(ctx *gin.Context, page *Page, canViewCreator bool, sub *model.CreatorSupport, session *sessions.VisitorSession, creator *model.User, campaign *model.Campaign) (data gin.H) {
	//load preview media if available
	var audio model.CreatorFile
	var recentSupporters []model.CreatorSupport

	//should just findone instead
	if files, err := campaign.Files(&(session.PrivateKey), canViewCreator); err == nil || len(files) != 0 {
		audio = files[0]
	} else {
		//file was deleted
		apiError(ctx, fmt.Sprintf("Internal server error. Error retrieving files for audio %s", campaign.ID))
	}
	streamURL, err := storage.GetVideoStreamURL(audio, canViewCreator)
	if err != nil {
		apiError(ctx, fmt.Sprintf("Internal server error. Error retrieving stream URL for audio %s", campaign.ID))
	}

	if !session.User.ID.IsZero() {
		session.User.LoggedIn = true

		recentSupporters, err = creator.ListRecentSupporters(8, 0, &(campaign.ID), true)
		if err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		}
	}

	data = gin.H{
		"audio":      audio,
		"stream_url": streamURL,
		"supporters": recentSupporters,
	}
	return

}

func photobookCampaign(ctx *gin.Context, page *Page, canViewCreator bool, sub *model.CreatorSupport, session *sessions.VisitorSession, creator *model.User, campaign *model.Campaign) (data gin.H) {
	var recentSupporters []model.CreatorSupport

	var photos []model.CreatorFile
	var err error
	if photos, err = campaign.Files(&(session.PrivateKey), canViewCreator); err != nil {
		apiError(ctx, "Internal server error. Error retrieving files")
		return
	}

	//only show only photo if not subscribed
	if !canViewCreator {
		photos = photos[:1]
	}

	if !session.User.ID.IsZero() {
		session.User.LoggedIn = true

		recentSupporters, err = creator.ListRecentSupporters(8, 0, &(campaign.ID), true)
		if err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		}
	}

	data = gin.H{
		"photos":     photos,
		"content":    campaign,
		"supporters": recentSupporters,
	}
	return
}

func embedCampaign(ctx *gin.Context, page *Page, canViewCreator bool, sub *model.CreatorSupport, session *sessions.VisitorSession, creator *model.User, campaign *model.Campaign) (data gin.H) {
	//load preview media if available

	var recentSupporters []model.CreatorSupport
	if !session.User.ID.IsZero() {
		session.User.LoggedIn = true

		var err error
		recentSupporters, err = creator.ListRecentSupporters(8, 0, &(campaign.ID), true)
		if err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		}
	}

	data = gin.H{
		"supporters": recentSupporters,
	}
	return
}

func UnlockContentForSession(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
	}

	supportID, _ := ctx.Params.Get("support_id")
	campaignID, _ := ctx.Params.Get("campaign_id")
	unlockCode := ctx.Param("unlock_code")
	username := ctx.Param("username")

	//signature := ctx.Query("signature")
	//ts := ctx.Query("ts")

	cid, _ := primitive.ObjectIDFromHex(campaignID)
	sid, _ := primitive.ObjectIDFromHex(supportID)
	creator, _ := model.RetrieveCreatorByUsername(username)

	supporter, _ := model.GetSupporter(sid, &(creator.ID))
	campaign, _ := model.CampaignByID(cid)

	if creator.ID.IsZero() || campaign.Owner != creator.ID {
		apiError(ctx, "User does not own this content")
		return
	}

	if supporter.UnlockCode.Hex() != unlockCode {
		apiError(ctx, "Invalid unlock code")
		return
	}

	if campaign.Subscription != "pay_per_view" && campaign.Subscription != "subscription" && campaign.Subscription != "service" {
		apiError(ctx, "Cannot unlock free content or service")
		return
	}

	err = session.UpdateGrantContentAccess(ctx, &cid)
	log.Printf("updated session with %v", err)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/@%s/%s", creator.Username, campaign.URI))
}
