package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

//PublicGetSiteConfig returns the site config
func PublicGetSiteConfig(ctx *gin.Context) {
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"site_mode":  "live",
		"public_url": "https://myhustlezw.co.zw/",
		"api_url":    "https://myhustlezw.co.zw/api/v1/",
		"site_name":  "MyHustle",
		"usd_to_zwl": cron.GetExchangeRate(),
	}, false)
}

func PublicSyncCurrency(ctx *gin.Context) {
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"usd_to_zwl": cron.GetExchangeRate(),
	}, false)
}

//RenderLandingPage renders the default home page
func RenderLandingPage(ctx *gin.Context) {
	page := newPage("Myhustle - The monetization platform for Zimbabwean content creators.")
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Page": page,
	})
}

//RenderAboutUs renders the about us page
func RenderAboutUs(ctx *gin.Context) {

}

//RenderCareers renders the careers page
func RenderCareers(ctx *gin.Context) {

}

//RenderExplore renders the explore creators page
func RenderExplore(ctx *gin.Context) {

}

//RenderContact renders the contact us page
func RenderContact(ctx *gin.Context) {

}

//RenderFAQ renders the frequently asked questions page
func RenderFAQ(ctx *gin.Context) {

}

//RenderDisputes renders the page for filing disputes
func RenderPaymentDisputes(ctx *gin.Context) {

}

//RenderCopyrightViolations renders the copyrights violations page
func RenderCopyrightViolations(ctx *gin.Context) {

}
