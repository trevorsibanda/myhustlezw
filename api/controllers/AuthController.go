package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"github.com/ttacon/libphonenumber"

	"github.com/gin-gonic/gin"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func apiError(ctx *gin.Context, reason string) {
	util.ApiError(ctx, reason)
}

//UserLoginForm models the login form
type UserLoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//AuthLoginView handles GET requests to login page
func AuthLoginView(ctx *gin.Context) {
	var flashes []interface{}

	vsession, err := sessions.GetVisitorSession(ctx)
	if err == nil {
		flashes = vsession.GinSession(ctx).Flashes()
	} else {
		vsession = sessions.NewVisitorSession(ctx)
	}

	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"flashes": flashes,
	})
}

//AuthSignupView handles GET requests to creator signup page
func AuthSignupView(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "creator-signup.html", gin.H{
		"user": nil,
	})
}

//AuthRecoverPassword renders the recover password page
func AuthRecoverPasswordView(ctx *gin.Context) {

}

//UserSignupForm models the details we collect when creating a fan account for a new user
type UserSignupForm struct {
	Fullname        string `form:"fullname" binding:"required"`
	Username        string `form:"username" binding:"required"`
	Email           string `form:"email" binding:"required"`
	CountryCode     string `form:"country_code" json:"country_code" binding:"required"`
	Role            string `form:"role" binding:"required"`
	PhoneNumber     string `form:"phone_number" json:"phone_number" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"password_confirm" json:"password_confirm" binding:"required"`
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func sanitizeSignupForm(form UserSignupForm) (newform UserSignupForm, err error) {
	form.Email = strings.TrimSpace(strings.ToLower(form.Email))
	form.Password = strings.TrimSpace(form.Password)
	form.Fullname = strings.TrimSpace(strings.Title(form.Fullname))

	if !isEmailValid(form.Email) {
		err = fmt.Errorf("%s is not a valid email address", form.Email)
		return
	}

	if len(form.Fullname) < 2 {
		err = fmt.Errorf("Fullname should be at least two characters long")
		return
	}

	if len(form.Password) < 6 {
		err = fmt.Errorf("Password should be at least 6 characters long")
		return
	}

	if form.Password != form.ConfirmPassword {
		err = fmt.Errorf("Your password do not match")
		return
	}

	num, err := libphonenumber.Parse(form.CountryCode+form.PhoneNumber, "ZW")
	if err != nil {
		err = fmt.Errorf("Phone number is not valid.")
	}
	form.PhoneNumber = fmt.Sprintf("+%d%d", num.GetCountryCode(), num.GetNationalNumber())

	if !model.CheckUsername(form.Username) {
		err = fmt.Errorf("Cannot use that username, either it has been registered or is not allowed.")

	}

	return form, err
}

//AuthSignupHandler creates new accounts for creators
func AuthSignupHandler(ctx *gin.Context) {
	var form UserSignupForm
	var sec *model.UserCredentials
	var creator *model.User
	var err error

	vsession := sessions.NewVisitorSession(ctx)
	vsession.LastActiveNow()

	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("bad form values: %v", err))
		return
	}
	if form, err = sanitizeSignupForm(form); err != nil {
		apiError(ctx, err.Error())
		return
	}

	if _, err = model.RetrieveUser(form.Email, form.Password, strings.ToLower(form.Username)); err == nil {
		apiError(ctx, "User registered with specified email or phone number or username already exists.")
		return
	}

	//create new user
	user := model.User{
		Fullname:         form.Fullname,
		PhoneNumber:      form.PhoneNumber,
		PhoneVerified:    false,
		Email:            form.Email,
		EmailVerified:    false,
		Created:          time.Now(),
		IdentityVerified: false,
		Profile:          model.DefaultCreatorProfile(form.Role),
		Page:             model.DefaultCreatorPage(),
		Subscriptions:    model.DefaultMembershipsConfig(),
		PayoutDetails: model.CreatorPayoutDetails{
			USD: model.CreatorPayoutFields{Method: model.PayoutBank},
			ZWL: model.CreatorPayoutFields{Method: model.PayoutBank},
		},
		AcceptingPayments: false,
	}
	user.Profile.Description = fmt.Sprintf("is a %s", form.Role)
	page := user.Page
	page.DonationItemName = "coffee"
	page.AllowSupporters = true
	page.SupportersName = "fan"
	page.ShowSocialMedia = true
	page.DonationItemUnitPrice = 1.00
	page.ThankYouMessage = fmt.Sprintf("Hey there \n Thank you for supporting my hustle.\n\nI apprecite your contribution and as a thank you(excclusive only to my true fans lol) you have access to this:\n\nFor my fans only:\n\nLINKS!!!")
	subs := user.Subscriptions
	subs.Price = 2.00
	subs.Period = 1
	subs.PeriodUnit = model.PeriodMonth
	subs.LastCount = 1
	subs.ThankYouMessage = page.ThankYouMessage

	user.Page = page
	user.Subscriptions = subs

	creator, sec, err = model.NewCreator(user, strings.ToLower(form.Username), form.Password)
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to store new user to database"))
		return
	}

	//authenticate user and create new session
	//create new session
	vsession.HasCreator = true
	vsession.User.LoggedIn = true
	vsession.User = *creator

	go func() {
		//send one time pins
		//todo: log this
		_, err := cron.SendEmailVerification(creator, model.AccountSignup)
		if err != nil {
			log.Printf("Failed to send email verification for %v . Reason: %v", creator, err)
			//just continue, but something must be done. Ideally update and refect we failed to send the email
		}

		_, err = cron.SendSMSVerification(creator, model.AccountSignup)
		if err != nil {
			log.Printf("Failed to send sms verification for %v . Reason: %v", creator, err)
			//just continue, but something must be done. Ideally update and refect we failed to send the sms
		}

		cron.PushNotifyNewCreatorLogin(*creator, vsession, sec)
	}()

	if vsession.Save(ctx) != nil {
		apiError(ctx, fmt.Sprintf("Failed to save visitor session"))
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, vsession.User, true)
	return
}

//AuthLogoutHandler destroys sessions and creates new ones
func AuthLogoutHandler(ctx *gin.Context) {
	vsession, err := sessions.GetVisitorSession(ctx)
	if err != nil {
		vsession.GinSession(ctx).Clear()
	}

	vsession = sessions.NewVisitorSession(ctx)
	vsession.Save(ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
	ctx.Abort()
}

type UserPasswordResetForm struct {
	Email string `json:"email" binding:"required"`
}

//AuthRequestPasswordReset sends a password reset email/sms
func AuthRequestPasswordReset(ctx *gin.Context) {
	var form UserPasswordResetForm
	var err error

	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("bad form values: %v", err))
		return
	}

	user, err := model.RetrieveUser(form.Email, "", "")
	if err != nil {
		apiError(ctx, fmt.Sprintf("User not found"))
		return
	}

	if ok, _ := model.CanSendVerification(&user, model.ResetPassword); ok {
		go func() {
			//send one time pins
			cron.SendEmailVerification(&user, model.ResetPassword)
		}()
	} else {
		apiError(ctx, fmt.Sprintf("You have exceeded the amount of verification requests. Please try again after two hours."))
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "ok"}, false)

}

type UserPasswordProcessResetForm struct {
	Email       string `json:"email" binding:"required"`
	Otp         string `json:"otp" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

//AuthProcessPasswordReset processes a password reset request
func AuthProcessPasswordReset(ctx *gin.Context) {
	var form UserPasswordProcessResetForm
	var err error

	if err = ctx.BindJSON(&form); err != nil {
		apiError(ctx, fmt.Sprintf("bad form values: %v", err))
		return
	}

	if len(form.NewPassword) < 6 {
		apiError(ctx, "Password should be at least 4 characters long")
		return
	}

	user, err := model.RetrieveUser(form.Email, "", "")
	if err != nil {
		apiError(ctx, fmt.Sprintf("User not found"))
		return
	}

	if otp, err := model.VerifyEmailCode(user.ID, user.Email, form.Otp, model.ResetPassword); err != nil {
		apiError(ctx, "Invalid verification code.")
		return
	} else {
		if !otp.Verified {
			apiError(ctx, "Verification code has expired.")
		} else {
			if err := user.UpdatePassword("*changedByPasswordReset*", form.NewPassword); err != nil {
				apiError(ctx, fmt.Sprintf("Failed to change password. %s", err))
				return
			}
			cron.SendResetPasswordNotification(user)
		}
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{"status": "changed"}, false)
}

//AuthLoginHandler handles login attempts and creates new sessions
func AuthLoginHandler(ctx *gin.Context) {
	var form UserLoginForm
	var creator model.User
	var err error

	var vsession sessions.VisitorSession
	if vsession, err = sessions.GetVisitorSession(ctx); err != nil {
		apiError(ctx, "Failed to retrieve user session")
		return
	}

	if ctx.BindJSON(&form) != nil {
		apiError(ctx, fmt.Sprintf("bad form values"))
		return
	}

	if creator, err = model.RetrieveUser(form.User, form.User, form.User); err != nil {
		apiError(ctx, "No account registered with specified email or phone number.")
		return
	}

	var cred model.UserCredentials
	if cred, err = model.RetrieveUserCredentials(creator.ID); err != nil {
		apiError(ctx, fmt.Sprintf("An internal error occured. Failed to fetch user credentials"))
		return
	}

	if model.ValidatePasswords(cred.Password, form.Password) {
		//log login time
		//vsession = sessions.NewVisitorSession(ctx)

		//create new session
		vsession.HasCreator = true
		vsession.User.LoggedIn = true
		vsession.User = creator
		vsession.UpdateInfo(ctx)

		//note this also saves the session
		vsession.UpdateGrantContentAccess(ctx, nil)

		go (func() {
			cron.PushNotifyNewCreatorLogin(creator, vsession, &cred)
		})()
		util.ScrubbedPublicAPIJSON(ctx, vsession.User, true)
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, gin.H{"error": "Incorrect password entered."}, false)
}
