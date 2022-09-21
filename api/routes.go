package api

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/gzip"
	"github.com/trevorsibanda/myhustlezw/api/controllers"
	"github.com/trevorsibanda/myhustlezw/api/controllers/api"
	"github.com/trevorsibanda/myhustlezw/api/middleware"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

var (
	siteMode = os.Getenv("MYHUSTLE_SITE_MODE")
)

//always keep sorted
var protectedCampaignRoutes = map[string]bool{
	"buymeacoffee":  true,
	"content":       true,
	"dispute":       true,
	"subscribe":     true,
	"subscriptions": true,
	"supporters":    true,
}

func Routes(apiVersion string, router *gin.Engine) {
	switch siteMode {
	case "maintanence":
		applyMaintanenceModeRoutes(router)
	case "comin_soon":
		applySiteComingSoonMode(router)
	default:
		applySiteRoutes(router)
		api := router.Group("api")
		{
			switch apiVersion {
			case "v1":
				applyV1Routes(api)
			default:
				log.Fatalf("Cannot apply routes for undefined version")
			}
		}

	}

}

func applyMaintanenceModeRoutes(router *gin.Engine) {

}

func applySiteComingSoonMode(router *gin.Engine) {

}

func cachedPage(ctx *gin.Context, expire time.Duration, handler gin.HandlerFunc) gin.HandlerFunc {
	return middleware.CacheForVisitors(ctx, expire, handler)
}

func applySiteRoutes(router *gin.Engine) {

	landing := router.Group("/").Use(gzip.Gzip(gzip.BestCompression))
	{
		landing.Any("/", serveWebsite("index.html"))
		landing.GET("/explore", serveWebsite("explore.html"))
		landing.GET("/about", serveWebsite("about.html"))
		landing.GET("/contact", serveWebsite("contact.html"))
		landing.GET("/faq", serveWebsite("faq.html"))
		landing.GET("/careers", serveWebsite("careers.html"))
		landing.GET("/login", middleware.RedirectIfLoggedIn, func(ctx *gin.Context) {
			ctx.Redirect(http.StatusMovedPermanently, "/auth/login")
		})
		landing.GET("/auth", serveReactApp)
		landing.GET("/auth/:path", serveReactApp)
		landing.GET("/logout", controllers.AuthLogoutHandler)
		landing.POST("/api/login", middleware.RedirectIfLoggedIn, controllers.AuthLoginHandler)
		landing.GET("/signup", middleware.RedirectIfLoggedIn, func(ctx *gin.Context) {
			ctx.Redirect(http.StatusMovedPermanently, "/auth/signup")
		})
		landing.POST("/api/signup", middleware.RedirectIfLoggedIn, controllers.AuthSignupHandler)

		landing.GET("/disputes/payment", serveWebsite("payment_disputes.html"))
		landing.GET("/disputes/copyright", serveWebsite("copyright_disputes.html"))
	}

	router.NoRoute(serveWebsite("404.html"))

	applyContentCreatorRoutes(router)
	applyUserPagesRoutes(router)
	applyStaticAssetsRoutes(router)
}

func applyUserPagesRoutes(router *gin.Engine) {

	viewer := router.Group("/v")
	viewer.Use(middleware.ActiveCreatorAccount)
	viewer.GET("/:username/:campaign_id/:support_id/:unlock_code", controllers.UnlockContentForSession)
	//todo: compression and generate html tags
	user := router.Group("@:username").Use(gzip.Gzip(gzip.BestCompression))
	user.Use(middleware.ActiveCreatorAccount)
	{
		user.GET("/:id", func(ctx *gin.Context) {
			uri := ctx.Param("id")
			var h gin.HandlerFunc
			if protectedCampaignRoutes[strings.ToLower(uri)] {
				h = cachedPage(ctx, time.Minute*5, controllers.PublicGetUser(controllers.GenerateCachedCreatorPage))
			} else {
				h = cachedPage(ctx, time.Minute*30, controllers.PublicGetCampaign(controllers.GenerateCachedCampaignPage))
			}

			h(ctx)
		})
		user.GET("", func(ctx *gin.Context) {
			h := cachedPage(ctx, time.Minute*5, controllers.PublicGetUser(controllers.GenerateCachedCreatorPage))
			h(ctx)
		})

	}
}

func serveReactApp(ctx *gin.Context) {
	//modify this to generate meta tags for user profile pages
	ctx.File(dashboardFile("index.html"))
}

func serveWebsite(file string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//modify this to generate meta tags for user profile pages
		ctx.File(siteFile(file))
	}
}

func applyContentCreatorRoutes(router *gin.Engine) {

	router.GET("/creator", middleware.AuthenticatedCreator, serveReactApp)
	creator := router.Group("/creator").Use(gzip.Gzip(gzip.BestCompression))
	creator.Use(middleware.AuthenticatedCreator)
	{
		creator.Any(":path", serveReactApp)
		creator.Any(":path/:ext", serveReactApp)
		creator.Any(":path/:ext/:another", serveReactApp)
		//router.GET("/creator/:path", serveReactApp)
	}

}

func applyStaticAssetsRoutes(router *gin.Engine) {
	router.Use(static.Serve("/assets/", static.LocalFile(assetFile(""), false)))
	router.Use(static.Serve("/static/", static.LocalFile(dashboardFile("static/"), false)))
	router.Use(static.Serve("/public/static/", static.LocalFile(siteFile("static/"), false)))
	router.GET("/stream/:acl/:campaign/:user/:content/:filename", controllers.StreamMedia)
	router.GET("/content/:acl/:user/:contentid/:tamper/:filename", controllers.ServeContent)
}

func dashboardFile(file string) string {
	return "./static/dashboard/build/" + file
}

func siteFile(file string) string {
	return "./static/website/" + file
}

func assetFile(file string) string {
	return dashboardFile("/assets/" + file)
}

func applyV1Routes(router *gin.RouterGroup) {
	private := router.Group("v1/private").Use(gzip.Gzip(gzip.BestCompression))
	router.POST("v1/user/web_subscription", middleware.AuthenticatedUser, api.UpdateUserWebSubscription)
	router.GET("v1/private/user", api.GetCurrentUser)
	//security
	router.GET("v1/private/security/phone/update/:phone", middleware.AuthenticatedUser, api.UpdatePhoneNumber)
	router.GET("v1/private/security/email/update/:email", middleware.AuthenticatedUser, api.UpdateEmail)
	router.GET("v1/private/security/phone/verify/:code", middleware.AuthenticatedUser, api.VerifySMSCode)
	router.GET("v1/private/security/email/verify/:code", middleware.AuthenticatedUser, api.VerifyEmailCode)
	router.GET("v1/private/security/phone/resend_sms", middleware.AuthenticatedUser, api.ResendSMSCode)
	router.GET("v1/private/security/email/resend_email", middleware.AuthenticatedUser, api.ResendEmailCode)
	router.POST("v1/private/security/password/update", middleware.AuthenticatedUser, api.ChangePassword)
	router.POST("v1/private/security/login", middleware.RedirectIfLoggedIn, controllers.AuthLoginHandler)
	router.POST("v1/private/security/request_reset_password", middleware.RedirectIfLoggedIn, controllers.AuthRequestPasswordReset)
	router.POST("v1/private/security/process_reset_password", middleware.RedirectIfLoggedIn, controllers.AuthProcessPasswordReset)
	router.GET("v1/private/security/identity/verify_by_payment/:phone", middleware.AuthenticatedUser, api.VerifyIdentityByPayment)
	router.POST("v1/private/security/signup", middleware.RedirectIfLoggedIn, controllers.AuthSignupHandler)

	//fans
	router.GET("v1/private/subscriptions/:filter", middleware.AuthenticatedUser, api.ListAllFanSubscriptions)
	router.GET("v1/private/subscription/:id", middleware.AuthenticatedUser, api.RetrieveSubscription)

	private.Use(middleware.AuthenticatedCreator)
	{
		//timeline
		private.GET("/feed/:view/:filter", controllers.FetchFilteredFeed)
		private.GET("/summary", api.AccountSummary)
		private.GET("/user/profile/set_avatar/:id", api.SetAvatar)
		//private.GET("/user", api.GetCurrentUser)
		private.POST("/user/page/basic", api.UpdateBasicUserPageDetails)
		private.GET("/user/page/publish/:bool", api.PublishPage)
		private.POST("/user/update_basics", api.UpdateUserBasics)
		private.POST("/user/update_page_configurables", api.UpdatePageConfigurables)
		private.POST("/user/update_notifications", api.UpdateUserNotifications)
		private.POST("/user/page/social_media", api.UpdatePageSocialMedia)
		private.POST("/user/page/layout", api.UpdatePageLayout)

		//supporters
		private.GET("/supporters/get/:id", api.GetSupporter)
		private.GET("/supporters/hide/:action/:id", api.HideSupporter)
		private.GET("/supporters/recent/:max_per_page", api.ListRecentSupporters)
		private.GET("/supporters/recent/:max_per_page/:skip", api.ListRecentSupporters)
		private.GET("/supporters/download_csv/:id", api.DownloadSupportersCSV)
		//services
		private.POST("/supporters/service/fulfill/:id", api.MarkServiceAsFulfilled)
		private.POST("/supporters/service/refund/:id", api.RefundService)

		//wallet
		private.GET("/wallet/summary", api.GetWalletSummary)
		private.GET("/wallet/operations", api.GetRecentWalletOperations)
		private.GET("/wallet/my_payments/:page_size", api.GetMyRecentPayments)
		private.GET("/wallet/operations/:page_size", api.GetRecentWalletOperations)
		private.POST("/wallet/withdrawal_request/:currency", api.RequestWithdrawal)
		private.GET("/wallet/payout_details", api.GetBankingWithdrawalDetails)
		private.POST("/wallet/payout_details/:currency", api.BankingWithdrawalDetails)
		private.POST("/wallet/activate_payments", api.ActivatePayments)
		//campaigns
		private.GET("/campaigns/all", api.ListCampaigns)
		private.GET("/campaign/supporters/:id/:max_per_page/:skip", api.ListRecentCampaignSupporters)
		private.POST("/campaign/update/:id", api.UpdateCampaign)
		private.GET("/campaign/delete/:id", api.DeleteCampaign)
		private.GET("/campaign/detailed/:id", api.GetDetailedCampaign)
		private.POST("/campaign/new/:type", api.CreateNewCampaign)

		//storage
		private.Any("/storage/upload/:type/:role", controllers.UploadToStorage)
		private.GET("/storage/upload_hook/:type/:provider/params", controllers.GenerateUploadURL)
		private.GET("/image/:width/:height/:id", controllers.AdaptiveImage)

		//config
		private.GET("/config/service_templates", api.ServiceTemplates)

	}

	//public resources dont need any authentication

	public := router.Group("v1/public").Use(gzip.Gzip(gzip.BestCompression))
	{

		public.GET("/image/:width/:height/:id", controllers.AdaptiveImage)
		public.GET("/stream/")
		public.POST("/campaign/contribute", func(ctx *gin.Context) {

		})
		public.POST("/transaction/initiate/:purpose", api.InitiateSupportPayment)
		public.GET("/transaction/xhr_poll/:id", api.PollPaymentStatus)
		public.POST("/mgodoyi/transaction/callback/ecocash", api.EcocashTransactionCallback)
		//buy a coffee without creating a new account. If authenticated asociate with new account
		public.POST("/user/buyacoffee", api.ListRecentSupporters)

		public.GET("/_admin/_/wallet/withdrawal/approve/:currency/:amount/:id", api.ApproveWithdrawal)

		//restricted content stream
		public.GET("/stream/public/:id.m3u8")
		public.GET("/stream/private/:id.m3u8")

		//public site api
		public.GET("/_site/config", controllers.PublicGetSiteConfig)
		public.GET("/_site/currency_sync", controllers.PublicSyncCurrency)
		public.GET("/_user/:username", middleware.ActiveCreatorAccount, controllers.PublicGetUser(util.ScrubbedPublicAPIJSON))
		public.GET("/_file/download/:username/:id", middleware.ActiveCreatorAccount, controllers.PublicDownloadCampaignFile)
		public.GET("/_campaign/:username/:id", middleware.ActiveCreatorAccount, controllers.PublicGetCampaign(util.ScrubbedPublicAPIJSON))
		//public.GET("/_gatekeeper/unlock_content/:payment_id/:campaign_id", controllers.UnlockContentForSession)
		public.GET("/_campaigns/:username/:page", middleware.ActiveCreatorAccount, controllers.PublicGetCampaign(util.ScrubbedPublicAPIJSON))
		public.POST("/_service/request/:username/:id", middleware.ActiveCreatorAccount, controllers.PublicRequestService)
		public.GET("/_service/request/:id", middleware.ActiveCreatorAccount, controllers.PublicGetServiceRequest)
		public.POST("/_campaign/metrics/:username/:id", middleware.ActiveCreatorAccount, controllers.PublicGetMetrics)
		public.POST("/_campaign/log_metric/:username/:id", middleware.ActiveCreatorAccount, controllers.PublicLogMetric)

	}

}
