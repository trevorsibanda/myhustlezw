package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListAllFanSubscriptions(ctx *gin.Context) {
	session := ctx.Keys["session"].(sessions.VisitorSession)

	filter := ctx.Param("filter")

	if subscriptions, err := session.User.GetSubscriptions(filter); err != nil {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to retrieve subscriptions. %v", err),
		}, false)
		return
	} else {
		var subs []gin.H
		for _, sub := range subscriptions {
			creator, _ := model.RetrieveCreatorByID(sub.Creator)
			subs = append(subs, gin.H{
				"sub": sub,
				"creator": gin.H{
					"username": creator.Username,
					"_id":      creator.ID,
					"profile":  creator.Profile,
					"type":     "creator",
				},
			})
		}
		util.ScrubbedPublicAPIJSON(ctx, subs, false)
	}
}

func RetrieveSubscription(ctx *gin.Context) {
	session := ctx.Keys["session"].(sessions.VisitorSession)
	var id primitive.ObjectID
	var err error

	idS, _ := ctx.Params.Get("id")
	if id, err = primitive.ObjectIDFromHex(idS); err != nil {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("invalid id passed."),
		}, true)
		return
	}

	if subscription, err := session.User.GetSubscription(id); err != nil {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to retrieve subscription. %v", err),
		}, true)
		return
	} else {
		creator, _ := model.RetrieveCreatorByID(subscription.Creator)
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"sub": subscription,
			"creator": gin.H{
				"username": creator.Username,
				"_id":      creator.ID,
				"profile":  creator.Profile,
				"type":     "creator",
			},
		}, false)
	}
}
