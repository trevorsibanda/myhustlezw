package middleware

import (
	"log"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/gin"
	siteCache "github.com/trevorsibanda/myhustlezw/api/cache"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
)

//CacheForVisitors caches a page for ordinary site visitors
func CacheForVisitors(ctx *gin.Context, expirate time.Duration, handler gin.HandlerFunc) (h gin.HandlerFunc) {

	session, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
		session.Save(ctx)
	}

	if session.User.ID.IsZero() {
		//cache the page
		log.Println("cached from here")
		h = cache.CachePage(siteCache.Store(), expirate, handler)
	} else {
		h = handler
	}
	return
}
