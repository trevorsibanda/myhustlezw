package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

var secureMiddleware *secure.Secure

func secureHTTPHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			c.Abort()
			return
		}

		// Avoid header rewrite if response is a redirection.
		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}

func InitSecureWebsite(router *gin.Engine) {

	secureMiddleware = secure.New(secure.Options{
		AllowedHosts:         []string{"myhustle\\.co\\.zw", ".*\\.myhustle\\.co\\.zw"},
		AllowedHostsAreRegex: true,
		HostsProxyHeaders:    []string{"X-Forwarded-Host"},
		SSLRedirect:          false,
		SSLHost:              "myhustle.co.zw",

		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,

		ContentSecurityPolicy: "-",
	})

	router.Use(secureHTTPHandler())

}
