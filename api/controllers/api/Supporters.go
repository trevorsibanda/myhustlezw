package api

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ListRecentSupporters list recent supporters up to :max_per_page
func ListRecentSupporters(ctx *gin.Context) {
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

	//todo add suppprt for cursors

	supporters, err := creator.ListRecentSupporters(maxItems, skip, nil, true)
	if err != nil {
		apiError(ctx, fmt.Sprintf("failed to list supporters"))
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"skip":       skip,
		"page_size":  maxItems,
		"supporters": supporters,
	}, true)

}

//MarkServiceAsFulfilled marks a service as fulfilled
func MarkServiceAsFulfilled(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)
	var supporterID primitive.ObjectID
	var supporter model.CreatorSupport
	var err error

	idS, _ := ctx.Params.Get("id")
	if supporterID, err = primitive.ObjectIDFromHex(idS); err != nil {
		apiError(ctx, "ID is invalid")
		return
	}

	if supporter, err = model.GetSupporter(supporterID, &(creator.ID)); err != nil {
		apiError(ctx, "Failed to find supporter")
		return
	}

	if supporter.Form == nil {
		apiError(ctx, "Type not supported for fulfillment")
		return
	}

	form := supporter.Form
	if form.Fulfilled {
		apiError(ctx, "Already fulfilled")
		return
	}

	if form.Refunded {
		apiError(ctx, "Cannot modify already refunded entry")
		return
	}

	form.Fulfilled = true
	form.Log = fmt.Sprintf("%s\nMarked as fulfilled %v", form.Log, session.Info.Dictionary())
	form.FulfilledAt = time.Now()
	supporter.LastUpdated = time.Now()
	supporter.Form = form

	if err := supporter.SaveChanges(); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save changes to database with error %v", err))
		return
	}

	go cron.OnServiceFulfilled(creator, supporter, session)
	util.ScrubbedPublicAPIJSON(ctx, supporter, true)
}

type refundServiceForm struct {
	Reason string `form:"reason" json:"reason" `
}

//RefundService refunds a support. Used for service offerings.
//Cancel then refund
func RefundService(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)
	var hform refundServiceForm
	var supporterID primitive.ObjectID
	var supporter model.CreatorSupport
	var err error

	ctx.BindJSON(&hform)

	idS, _ := ctx.Params.Get("id")
	if supporterID, err = primitive.ObjectIDFromHex(idS); err != nil {
		apiError(ctx, "ID is invalid")
		return
	}

	if supporter, err = model.GetSupporter(supporterID, &(creator.ID)); err != nil {
		apiError(ctx, "Failed to find supporter")
		return
	}

	if supporter.Form == nil {
		apiError(ctx, "Type not supported for fulfillment")
		return
	}

	form := supporter.Form
	if form.Fulfilled {
		apiError(ctx, "Already fulfilled")
		return
	}

	if form.Refunded {
		apiError(ctx, "Cannot modify already refunded entry")
		return
	}

	form.Refunded = true
	form.Log = fmt.Sprintf("%s\nMarked as refunded %v", form.Log, session.Info.Dictionary())
	form.RefundedAt = time.Now()
	supporter.LastUpdated = time.Now()
	supporter.Form = form

	//mark the wallet op as refunded
	walletOp, err := creator.GetWalletOp(supporter.WalletOp)
	if err != nil {
		apiError(ctx, "Failed to refund because failed to rerieve wallet op")
		return
	}

	walletOp.Dispute("Cancelled by creator")
	if err = walletOp.ApproveRefund(hform.Reason); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to mark wallet op as refunded with error %v", err))
		return
	}

	if err := supporter.SaveChanges(); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save changes to database with error %v", err))
		return
	}

	go cron.OnServiceRefunded(creator, supporter, session)
	util.ScrubbedPublicAPIJSON(ctx, supporter, true)
}

func GetSupporter(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	var supporterID primitive.ObjectID
	var supporter model.CreatorSupport
	var err error

	idS, _ := ctx.Params.Get("id")
	if supporterID, err = primitive.ObjectIDFromHex(idS); err != nil {
		apiError(ctx, "ID is invalid")
		return
	}

	if supporter, err = model.GetSupporter(supporterID, &(creator.ID)); err != nil {
		apiError(ctx, "Failed to find supporter")
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, supporter, true)
}

func HideSupporter(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	var supporterID primitive.ObjectID
	var supporter model.CreatorSupport
	var err error

	idS, _ := ctx.Params.Get("id")
	if supporterID, err = primitive.ObjectIDFromHex(idS); err != nil {
		apiError(ctx, "ID is invalid")
		return
	}

	action := ctx.Param("action")

	if supporter, err = model.GetSupporter(supporterID, &(creator.ID)); err != nil {
		apiError(ctx, "Failed to find supporter")
		return
	}

	if supporter.Creator != creator.ID {
		apiError(ctx, "You do not have permission to do this")
		return
	}
	switch action {
	case "message":
		supporter.HideMessage = !supporter.HideMessage
	case "all":
		supporter.Hidden = !supporter.Hidden
	default:
		apiError(ctx, "Unknown action")
		return
	}

	supporter.SaveChanges()
	util.ScrubbedPublicAPIJSON(ctx, supporter, true)
}

//DownloadSupportersCSV generates a csv of all or campaign specific supporters
func DownloadSupportersCSV(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	var supporters []model.CreatorSupport
	var campaign model.Campaign
	var err error

	if id, _ := ctx.Params.Get("id"); id == "all" {
		supporters, err = creator.ListRecentSupporters(math.MaxInt32, 0, nil, true)
	} else {
		var cID primitive.ObjectID

		if cID, err = primitive.ObjectIDFromHex(id); err == nil {
			apiError(ctx, fmt.Sprintf("invalid campaign id"))
			return
		}
		if campaign, err = model.CampaignByID(cID); err != nil {
			apiError(ctx, fmt.Sprintf("resourcse does not exist"))
			return
		}
		if campaign.Owner != creator.ID {
			apiError(ctx, fmt.Sprintf("permission denied"))
			return
		}
		supporters, err = campaign.ListRecentSupporters(math.MaxInt64, 0)
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	for _, supporter := range supporters {
		row := []string{
			supporter.Type,
			supporter.DisplayName,
			supporter.Email,
			supporter.Comment,
			supporter.Currency,
			fmt.Sprintf("%.2f", supporter.Amount),
			supporter.Created.String(),
		}
		writer.Write(row)
	}

	filename := fmt.Sprintf("supporters__%s_%s.csv", campaign.URI, time.Now().String())
	writer.Flush()
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "text/csv", buffer.Bytes())

}
