package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetWalletSummary returns the wallet summary for all defined currencies
func GetWalletSummary(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	data := gin.H{}
	for _, currency := range []string{"USD", "ZWL"} {
		summary, err := creator.WalletSummary(currency)
		if err != nil {
			//log the error here
			continue
		}
		data[currency] = summary
	}

	util.ScrubbedPublicAPIJSON(ctx, data, true)
}

//GetRecentWalletOperations returns the recent wallet operations
func GetRecentWalletOperations(ctx *gin.Context) {
	recentOps := make([]model.WalletOperation, 0)
	var err error
	creator := ctx.Keys["creator"].(model.User)
	var pageSize int = 10

	if maxPage, ok := ctx.Params.Get("page_size"); ok {
		pageSize, _ = strconv.Atoi(maxPage)
	}

	recentOps, err = creator.ListRecentWalletOperations(int64(pageSize))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, recentOps, true)
}

//GetMyRecentPayments returns the recent payments for current user
func GetMyRecentPayments(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	var pageSize int = 10

	if maxPage, ok := ctx.Params.Get("page_size"); ok {
		pageSize, _ = strconv.Atoi(maxPage)
	}

	recentPayments, err := creator.ListRecentPayments(int64(pageSize))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, recentPayments, true)
}

func getMinWithdrawal(currency model.PaymentCurrency) float64 {
	switch currency {
	case model.USD:
		return 20.00 //20usd
	case model.ZWL:
		return cron.USDToZWL(5.00) //min withdrawal is $5
	default:
		return 1000000
	}
}

//RequestWithdrawal requests to widraw all funds in the account
func RequestWithdrawal(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	var currency model.PaymentCurrency
	var summary model.CreatorWalletSummary
	var err error

	if currencyStr, ok := ctx.Params.Get("currency"); ok {
		switch strings.ToLower(currencyStr) {
		case "usd":
			currency = model.USD
		case "zwl":
			currency = model.ZWL
		default:
			apiError(ctx, fmt.Sprintf("Invalid currency %s", currencyStr))
			return
		}
	} else {
		apiError(ctx, fmt.Sprintf("No currency supplied %s", err))
		return
	}

	if summary, err = creator.WalletSummary(string(currency)); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to fetch wallet for currency %s", err))
		return
	}

	if summary.PendingWithdrawal > 0.00 && summary.Available != 0.00 {
		apiError(ctx, "There is already a withdrawal pending. Cannot create a new one.")
		return
	}

	if min := getMinWithdrawal(currency); summary.Available <= min {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{"error": fmt.Sprintf("You have %s. Min withdrawal amount is %s", cron.Format(string(currency), summary.Available), cron.Format(string(currency), min))}, true)
		return
	}

	//process withdrawal
	err = model.RequestWithdrawal(creator, currency, "Creator requested withdrawal from web api")
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to process withdrawal. Reason: %v", err))
		return
	}

	summary, _ = creator.WalletSummary(string(currency))

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status":  "ok",
		"summary": summary,
	}, true)

	go cron.SendWalletWithdrawalRequest(creator, summary)

}

type bankingDetailsForm struct {
	Method          string `form:"method" binding:"required"`
	BankName        string `form:"bankname" `
	BankBranch      string `form:"bankbranch" `
	BankAccount     string `form:"bankaccountnumber" json:"bankaccountnumber"`
	BankAccountName string `form:"bankaccountname" `
	PhoneNumber     string `form:"phonenumber" `
}

//BankingWithdrawalDetails updates cashout details for a given currency
func BankingWithdrawalDetails(ctx *gin.Context) {
	var form bankingDetailsForm

	creator := ctx.Keys["creator"].(model.User)

	var currency model.PaymentCurrency
	if currencyStr, ok := ctx.Params.Get("currency"); !ok {
		apiError(ctx, fmt.Sprintf("currency not provided"))
		return
	} else {
		if currencyStr == "USD" {
			currency = model.USD
		} else {
			currency = model.ZWL
		}
	}

	if err := ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("invalid form %v", err))
		return
	}

	var method model.PayoutMethod
	if form.Method == "bank" {
		method = model.PayoutBank
	} else if form.Method == "mobile_money" {
		method = model.PayoutMobileMoney
	} else {
		apiError(ctx, fmt.Sprintf("invalid payout method %s", form.Method))
		return
	}

	details := model.CreatorPayoutFields{
		Method:            method,
		BankName:          form.BankName,
		BankBranch:        form.BankBranch,
		BankAccountName:   form.BankAccountName,
		BankAccountNumber: form.BankAccount,
		PhoneNumber:       form.PhoneNumber,
	}

	if err := creator.UpdatePayoutDetails(currency, details); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to update payout details %v", err))
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "ok",
	}, true)

	go cron.SendPayoutDetailsChangedNotification(creator, currency, details)
}

//GetBankingWithdrawalDetails gets the payout details for a user
func GetBankingWithdrawalDetails(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"USD": creator.PayoutDetails.USD,
		"ZWL": creator.PayoutDetails.ZWL,
	}, true)
}

func ApproveWithdrawal(ctx *gin.Context) {
	var creatorID primitive.ObjectID
	var creator model.User
	var currency model.PaymentCurrency
	var amount float64
	var err error
	var summary model.CreatorWalletSummary

	if creatorID, err = primitive.ObjectIDFromHex(ctx.Param("id")); err != nil {
		apiError(ctx, "Invalid ID")
		return
	}

	if creator, err = model.RetrieveCreatorByID(creatorID); err != nil {
		apiError(ctx, "Failed to fetch user")
		return
	}

	switch strings.ToLower(ctx.Param("currency")) {
	case "zwl":
		currency = model.ZWL
	case "usd":
		currency = model.USD
	default:
		apiError(ctx, "Invalid currency")
		return
	}

	if amount, err = strconv.ParseFloat(ctx.Param("amount"), 32); err != nil {
		apiError(ctx, "Invalid amount")
		return
	}

	summary, err = creator.WalletSummary(string(currency))
	if err != nil {
		apiError(ctx, "Failed to get wallet summary")
		return
	}

	//now process the withdrawals
	var count int64
	if count, err = summary.ApprovePendingWithdrawal(creatorID); err != nil {
		apiError(ctx, "Request was not successful. Withdrawal not processed!!!")
		return
	}

	cron.SendWalletWithdrawalApproved(creator, summary, fmt.Sprintf("%d transactions approved for withdrawal", count))
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status":  "ok",
		"message": fmt.Sprintf("%s marked as withdrawn from @%s account. Request had %f in url. Creator and finance dpt notified", cron.Format(summary.Currency, summary.PendingWithdrawal), creator.Username, amount),
	}, true)

}

//ActivatePayments activates payments for a user account. This includes verifying the email address
func ActivatePayments(ctx *gin.Context) {

}
