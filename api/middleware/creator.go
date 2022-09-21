package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

//AuthenticatedCreator is a middleware to allow only authenticated creators
//access to a resource
func AuthenticatedCreator(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
		//ctx.Redirect(http.StatusMovedPermanently, "/login?no_session")
		//ctx.Abort()
		//return
	}

	if session.User.ID.IsZero() {
		target := "/login"

		if strings.HasPrefix(ctx.Request.RequestURI, "/api/") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, target)
		ctx.Abort()
		return
	}
	creator, err := model.RetrieveCreatorByID(session.User.ID)

	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login?user_deleted")
		return
	} else {
		session.User.LoggedIn = true
	}
	ctx.Keys["creator"] = creator
	session.User = creator
	session.Save(ctx)
	ctx.Keys["session"] = session
	ctx.Next()
}

//AuthenticatedUser is a middleware to only allow access to authenticated users.
//this includes fans and creators
func AuthenticatedUser(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
		//ctx.Redirect(http.StatusMovedPermanently, "/login?no_session")
		//ctx.Abort()
		//return
	}

	if session.User.ID.IsZero() {
		if strings.HasPrefix(ctx.Request.RequestURI, "/api/v1/private/") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Not authenticated"})
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		ctx.Abort()
		return
	} else {
		session.User.LoggedIn = true
	}

	ctx.Keys["session"] = session
	ctx.Next()
}

//RedirectIfLoggedIn is a middleware to redirect users to their dashboard
func RedirectIfLoggedIn(ctx *gin.Context) {
	var dashboardURL string = "/creator/dashboard"

	session, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
		ctx.Next()
		return
	}

	if session.User.ID.IsZero() {
		ctx.Next()
		return
	}

	_, err = model.RetrieveCreatorByID(session.User.ID)
	if err == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, dashboardURL)
		return
	}

	ctx.Next()
}

func apiError(ctx *gin.Context, reason interface{}) {
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"error": reason,
	}, false)
	return
}

//ActiveCreatorAccount middleware limits access of a creator's page when the account has
//not yet beign verified, not published, or otherwise
func ActiveCreatorAccount(ctx *gin.Context) {
	session, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		session = sessions.NewVisitorSession(ctx)
		session.Save(ctx)
	}
	username, ok := ctx.Params.Get("username")
	if !ok {
		apiError(ctx, "Error in implementation of routes! username should be available for all public creator routes")
		return
	}

	creator, err := model.RetrieveCreatorByUsername(username)

	if err != nil {
		apiError(ctx, fmt.Sprintf("User %s does not exist", username))
		return
	}
	ctx.Keys["creator"] = creator
	ctx.Keys["session"] = session
	ctx.Next()
}
