package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchFilteredFeed(ctx *gin.Context) {
	session, _ := sessions.GetVisitorSession(ctx)
	filter := model.FilterFromString(ctx.Param("filter"))
	view := model.ViewFromString(ctx.Param("view"))
	creatorFilterStr := ctx.Query("creator")
	var targetCreator *primitive.ObjectID
	if filter == model.FeedFilterByCreator && view == model.FeedViewFiltered {
		tc, err := model.RetrieveCreatorByUsername(creatorFilterStr)
		if err != nil {
			apiError(ctx, fmt.Sprintf("User %s does not exist", creatorFilterStr))
			return
		}
		if tc.ID.IsZero() {
			apiError(ctx, "User does not exist. Null ID")
			return
		}
		targetCreator = &(tc.ID)
	}

	skip, _ := strconv.Atoi(ctx.Param("skip"))

	session.UpdateGrantContentAccess(ctx, nil)
	paidContent := session.PaidContent
	subs, _ := session.User.GetSubscriptions("active")
	subscribedUsers := make([]primitive.ObjectID, 0)

	for _, s := range subs {
		subscribedUsers = append(subscribedUsers, s.Creator)
	}

	//TODO: store users interacted with in session
	interactedUsers := make([]primitive.ObjectID, 0)

	rf, err := session.User.GenerateUserFeed(int64(skip), view, filter, paidContent, subscribedUsers, interactedUsers, targetCreator)
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to generate feed with error %v", err))
		return
	} else {
		for idx, rfi := range rf {

			canViewCreator := false
			for _, s := range subs {
				if s.Creator == rfi.Content.Owner {
					canViewCreator = true
					break
				}
			}
			log.Printf("canViewCreator %v %v", subs, canViewCreator)
			canView, _ := CanViewCreatorSubscriberContent(rfi.Creator, session, &(rfi.Content), canViewCreator)
			rfi.Content.CanView = canView
			rfi.Content.PreviewImageURL = rfi.Content.PreviewURL(&session.PrivateKey, canView, 480, 480)
			rf[idx] = rfi
		}
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"feed":    rf,
		"view":    ctx.Param("view"),
		"filter":  ctx.Param("filter"),
		"skip":    skip,
		"creator": creatorFilterStr,
	}, false)

}
