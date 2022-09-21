package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/storage"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

//CreatorPageView renders the creator's profile page
func PublicGetUser(writer serveWriter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creator := ctx.Keys["creator"].(model.User)
		session := ctx.Keys["session"].(sessions.VisitorSession)
		//todo load other necessary content

		var recentSupporters []model.CreatorSupport
		var recentCampaigns []model.Campaign
		var featured, subsheadline model.CreatorFile
		var err error
		var streamURL string

		page := newPage(fmt.Sprintf("@%s ", creator.Username))

		hasSubscription, sub := canViewCreatorSubscriberContent(creator, session, nil, false)

		if !creator.Page.FeaturedContent.IsZero() {
			featured, err = model.GetFileByID(creator.Page.FeaturedContent, nil, true)
			if err == nil && featured.Type == "video" {
				streamURL, _ = storage.GetVideoStreamURL(featured, true)
				page.AddPrefetch(streamURL)
			}
			page.AddPrefetch(featured.URL)
		}

		if creator.Subscriptions.Active {
			subsheadline, err = model.GetFileByID(creator.Subscriptions.HeadlineImage, nil, true)
			creator.Subscriptions.HeadlineImageURL = subsheadline.URL
			page.AddPrefetch(subsheadline.URL)
		} else {
			subsheadline = model.CreatorFile{
				Type:    "image",
				Storage: "local",
				ID:      creator.Profile.ProfilePhoto,
				Owner:   creator.ID,
			}
		}

		if !session.User.ID.IsZero() {
			session.User.LoggedIn = true
		}

		if recentSupporters, err = creator.ListRecentSupporters(8, 0, nil, session.User.LoggedIn); err != nil {
			fmt.Printf("Failed to load recent supporters for %v. Error %v", creator, err)
		} else {
			for idx, support := range recentSupporters {
				if support.HideMessage {
					support.Comment = ""
					recentSupporters[idx] = support
				}
			}
		}

		recentCampaigns, _ = creator.ListRecentCampaigns(&(session.PrivateKey), 10, 0, hasSubscription)

		//TODO: Perfomance hit by checking database every time, use memcached
		if !session.User.ID.IsZero() {
			session.UpdateGrantContentAccess(ctx, nil)
		}

		for idx, c := range recentCampaigns {
			//update if the campaign can be viewed individually
			c.CanView, _ = CanViewCreatorSubscriberContent(creator, session, &c, hasSubscription)
			c.PreviewImageURL = c.PreviewURL(&(session.PrivateKey), c.CanView, 480, 480)
			page.AddPrefetch(c.PreviewImageURL)
			recentCampaigns[idx] = c
		}

		page.SetPageAuthorURL(creator.URL())
		page.SetImage(creator.ProfilePhotoURL())
		page.AddPrefetch(creator.ProfilePhotoURL())
		page.SetImageAlt(creator.Fullname)
		page.SetPageURL(creator.URL())
		page.SetTwitterUsername(creator.Profile.SocialMedia.Twitter)
		page.SetDescription(creator.Profile.AboutMe)

		data := gin.H{
			"page":          &page,
			"creator":       util.ScrubPublic(creator),
			"subscription":  sub,
			"featured":      featured,
			"subs_headline": subsheadline,
			"stream_url":    streamURL,
			"campaigns":     recentCampaigns,
			"supporters":    recentSupporters,
			"user":          session.User,
			"allowed":       hasSubscription,
		}
		writer(ctx, data, session.User.LoggedIn)
	}
}

func GenerateCachedCreatorPage(ctx *gin.Context, data interface{}, loggedIn bool) {
	file, err := os.ReadFile(dashboardFile("index.html"))
	if err != nil {
		apiError(ctx, "Failed to render page.")
		return
	}
	page := data.(gin.H)["page"].(*Page)
	meta := util.GenerateMeta(page.Dict)
	var head strings.Builder
	head.WriteString(meta)
	for _, link := range page.Prefetch {
		head.WriteString(fmt.Sprintf("<link rel=\"prefetch\" href=\"%s\">\n", link))
	}

	editted := strings.Replace(string(file), "<pagemeta></pagemeta>", head.String(), 1)
	withPreGen := util.GeneratePageWithData("creator", editted, data)

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(withPreGen))
}

//CreatorSupportPage shows the buy me a coffee page
func CreatorSupportPage(ctx *gin.Context) {

}

//CreatorAboutMe shows the creators about me page
func CreatorAboutMe(ctx *gin.Context) {

}

//CreatorThankYou shows the creator's thank you page after the user pays for support or subscribes
func CreatorThankYou(ctx *gin.Context) {

}

//CreatorAllContentPage lists all the creator's content in one page
func CreatorAllContentPage(ctx *gin.Context) {

}

func canViewCreatorSubscriberContent(creator model.User, session sessions.VisitorSession, campaign *model.Campaign, hasSubscription bool) (b bool, sub model.CreatorSupport) {
	return CanViewCreatorSubscriberContent(creator, session, campaign, hasSubscription)
}

func CanViewCreatorSubscriberContent(creator model.User, session sessions.VisitorSession, campaign *model.Campaign, hasSubscription bool) (b bool, sub model.CreatorSupport) {
	if creator.ID == session.User.ID {
		b = true
		return
	}

	if campaign != nil {
		lookupMap := make(map[string]bool)
		for _, id := range session.PaidContent {
			lookupMap[id.Hex()] = true
		}
		if exists, _ := lookupMap[campaign.ID.Hex()]; exists {
			campaign.CanView = true
			campaign.PreviewImageURL = campaign.PreviewURL(&(session.PrivateKey), true, 480, 480)
			b = true
			return
		}

		if campaign.Subscription == "public" {
			b = true
			return
		}
		if campaign.Subscription == "fans" && hasSubscription && !session.User.ID.IsZero() {
			b = true
		}
		return
	}

	//anonymous and not paid to view
	if session.User.ID.IsZero() && !b {
		b = false
		return
	} else if !session.User.ID.IsZero() && creator.ID == session.User.ID { //not anon and is creator
		b = true
		return
	}

	//if logged in and user has subscription which has not yet expired
	// and campaign is not pay per view
	if !b && !session.User.ID.IsZero() {
		sub, _ = session.User.GetSubscription(creator.ID)
		b = (!sub.ID.IsZero()) && time.Now().Before(sub.Expires)
		if b && campaign != nil {
			b = b && campaign.Subscription != "pay_per_view"
		}
	}

	return
}

func isPageOwner(creator *model.User, session *sessions.VisitorSession) bool {
	return session.User.ID == creator.ID
}

func siteFile(file string) string {
	return "./static/dashboard/build/" + file
}
