package api

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"github.com/ttacon/libphonenumber"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//UpdatePhoneNumber updates the phone number and sends a OTP
func UpdatePhoneNumber(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to update your phone number")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	phoneNumber := ctx.Param("phone")

	num, err := libphonenumber.Parse(phoneNumber, "ZW")
	if err != nil {
		err = fmt.Errorf("Phone number is not valid. Only Zimbabwean phone numbers allowed")
		apiError(ctx, err)
		return
	}
	phoneNumber = fmt.Sprintf("+%d%d", num.GetCountryCode(), num.GetNationalNumber())

	if _, err := model.RetrieveUser("", phoneNumber, ""); err != nil {
		apiError(ctx, fmt.Sprintf("A user with that phone number already exists"))
		return
	}

	if ok, _ := model.CanSendVerification(&user, model.ChangePhone); !ok {
		apiError(ctx, fmt.Sprintf("You have exceeded the amount of verification requests. Please try again after two hours."))
		return
	}

	err = user.UpdatePhoneNumber(phoneNumber)
	if err != nil {
		apiError(ctx, err)
		return
	}
	user.PhoneNumber = phoneNumber
	action := model.AccountSignup
	if user.PhoneVerified {
		action = model.ChangePhone
		cron.SendPhoneNumberChangedEmail(user)
	}

	_, err = cron.SendSMSVerification(&user, action)
	if err != nil {
		log.Printf("Failed to send sms verification for %v . Reason: %v", user, err)
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to send SMS but updated phoneNumber from %s to %s", user.PhoneNumber, phoneNumber),
		}, true)
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "updated",
	}, true)

}

//UpdateEmail updates the email and sends a OTP
func UpdateEmail(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	if ok, _ := model.CanSendVerification(&user, model.ChangeEmail); !ok {
		apiError(ctx, fmt.Sprintf("You have exceeded the amount of verification requests. Please try again after two hours."))
		return
	}

	email := strings.TrimSpace(strings.ToLower(ctx.Param("email")))

	if user, err := model.RetrieveUser(email, "", ""); err != nil && email != user.Email {
		apiError(ctx, fmt.Sprintf("A user with that email already exists"))
		return
	}

	if !isEmailValid(email) {
		apiError(ctx, fmt.Sprintf("%s is an invalid email", email))
		return
	}

	err = user.UpdateEmail(email)
	if err != nil {
		apiError(ctx, err)
		return
	}
	oldEmail := user.Email
	user.Email = email
	action := model.AccountSignup
	if user.EmailVerified {
		action = model.ChangePhone
		cron.SendEmailChangedEmail(oldEmail, user)
	}

	_, err = cron.SendEmailVerification(&user, action)
	if err != nil {
		log.Printf("Failed to send email verification for %v . Reason: %v", user, err)
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to send email but updated from %s to %s", user.Email, email),
		}, true)
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "updated",
	}, true)
}

//VerifyIdentityByPayment verifies the identity of the user by payment
func VerifyIdentityByPayment(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	if user.IdentityVerified {
		apiError(ctx, "You have already verified your identity")
		return
	}

	phone := ctx.Param("phone")
	if strings.HasPrefix(phone, "+") {
		if !strings.HasPrefix(phone, "+263") {
			apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", phone))
			return
		}
	}
	if strings.HasPrefix(phone, "0") {
		if strings.HasPrefix(phone, "0263") || strings.HasPrefix(phone, "00263") {
			apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", phone))
			return
		}
		if !strings.HasPrefix(phone, "078") && !strings.HasPrefix(phone, "077") {
			apiError(ctx, fmt.Sprintf("Sorry only Zimbabwean numbers allowed for Ecocash payments. %s", phone))
			return
		}
		phone = "+263" + phone[1:]
	}

	var lphone *libphonenumber.PhoneNumber
	if lphone, err = libphonenumber.Parse(phone, "ZW"); err != nil {
		apiError(ctx, fmt.Sprintf("Invalid phone number passed"))
		return
	}

	var paymentID primitive.ObjectID
	expires := time.Now().Add(time.Minute * 5)
	payment := model.PendingPayment{
		Created:             time.Now(),
		Comment:             "Verify Identity",
		Gateway:             model.GatewayEcocash,
		Email:               user.Email,
		Fullname:            user.Fullname,
		PhoneCountryCode:    fmt.Sprintf("+%d", lphone.GetCountryCode()),
		PhoneNationalNumber: fmt.Sprintf("%d", lphone.GetNationalNumber()),
		Currency:            model.ZWL,
		Price:               cron.USDToZWL(1),
		TargetFan:           user.ID,
		Action:              model.VerifyAccountOnPaid,
		Status:              "queued",
	}

	if paymentID, err = model.NewVerifyAccountPendingPayment(payment); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to save new payment"))
		return
	}
	payment.ID = paymentID

	md5Sum := md5.Sum([]byte(fmt.Sprintf("%snonce%d%s%d", paymentID.Hex(), expires.Unix(), paymentID.Hex(), expires.Unix())))

	go cron.DispatchPayment(payment, session, user)

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"_id":       paymentID,
		"gateway":   payment.Gateway,
		"amount":    payment.Price,
		"status":    payment.Status,
		"ts":        expires.Unix(),
		"signature": fmt.Sprintf("%x", md5Sum),
	}, true)
}

//VerifySMSCode verifies the code received via SMS
func VerifySMSCode(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	verificationCode := ctx.Param("code")

	var resp gin.H

	action := model.AccountSignup
	if user.PreviouslyVerified {
		action = model.ChangePhone
	}

	if otp, err := model.VerifySMSCode(user.ID, user.PhoneNumber, verificationCode, action); err != nil {
		resp = gin.H{
			"error": "Invalid verification code.",
		}
	} else {
		if !otp.Verified {
			resp = gin.H{"error": "Verification code has expired."}
		} else {
			resp = gin.H{"status": "verified"}
			log.Println(otp.Action)
			switch otp.Action {
			case model.AccountSignup:
				err = cron.OnSignupVerifiedPhone(user, session)
			case model.ChangePhone:
				cron.OnPhoneChanged(user)
			case model.ResetPassword:
				cron.OnResetPasswordPhone(user)
			case model.ChangeEmail:
				cron.OnEmailChanged(user)
			}
		}
	}

	util.ScrubbedPublicAPIJSON(ctx, resp, true)
}

//VerifyEmailCode verifies the code received via SMS
//todo: add rate limitting
func VerifyEmailCode(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	verificationCode := ctx.Param("code")

	action := model.AccountSignup
	if user.PreviouslyVerified {
		action = model.ChangeEmail
	}

	var resp gin.H

	if otp, err := model.VerifyEmailCode(user.ID, user.Email, verificationCode, action); err != nil {
		resp = gin.H{
			"error": "Invalid verification code.",
		}
	} else {
		if !otp.Verified {
			resp = gin.H{"error": "Verification code has expired."}
		} else {
			resp = gin.H{"status": "verified"}
			switch otp.Action {
			case model.AccountSignup:
				err = cron.OnEmailVerified(user)
			case model.ChangePhone:
				cron.OnPhoneChanged(user)
			case model.ResetPassword:
				cron.OnResetPasswordPhone(user)
			case model.ChangeEmail:
				cron.OnEmailChanged(user)
			}
		}
	}

	util.ScrubbedPublicAPIJSON(ctx, resp, true)

}

//ResendSMSCode resends the sms code for the current number
//todo: rate limit this
func ResendSMSCode(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	action := model.AccountSignup
	if user.PreviouslyVerified {
		action = model.ChangePhone
	}

	if ok, _ := model.CanSendVerification(&user, model.ChangePhone); !ok {
		apiError(ctx, fmt.Sprintf("You have exceeded the amount of verification requests. Please try again after two hours."))
		return
	}

	if user.PreviouslyVerified {
		go cron.SendPhoneNumberChangedEmail(user)
	}

	if _, err = cron.SendSMSVerification(&user, action); err != nil {
		log.Printf("Failed to send sms verification for %v . Reason: %v", user, err)
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to send SMS but  to %s", user.PhoneNumber),
		}, true)
		return
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "sent",
	}, true)
}

//ResendEmailCode resends the email code for the current email
//todo: rate limit this
func ResendEmailCode(ctx *gin.Context) {
	var user model.User
	var err error

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	action := model.AccountSignup
	if user.PreviouslyVerified {
		action = model.ChangePhone
	}

	if ok, _ := model.CanSendVerification(&user, action); !ok {
		apiError(ctx, fmt.Sprintf("You have exceeded the amount of verification requests. Please try again after two hours."))
		return
	}

	if _, err = cron.SendEmailVerification(&user, action); err != nil {
		log.Printf("Failed to send email verification for %v . Reason: %v", user, err)
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to send email but changed email to %s", user.Email),
		}, true)
		return
	}

	if user.PreviouslyVerified {
		go cron.SendEmailChangedEmail(user.Email, user)
	}

	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "sent",
	}, true)
}

type changePasswordForm struct {
	Password    string `form:"old" json:"old" binding:"required"`
	NewPassword string `form:"new" json:"new" binding:"required"`
}

//ChangePassword changes the user's password
func ChangePassword(ctx *gin.Context) {
	var user model.User
	var err error

	var form changePasswordForm
	if err = ctx.BindJSON(&form); err != nil {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Invalid form"),
		}, true)
		return
	}

	session := ctx.Keys["session"].(sessions.VisitorSession)
	if session.User.ID.IsZero() {
		apiError(ctx, "You are not allowed to access this resource")
		return
	} else if user, err = model.RetrieveCreatorByID(session.User.ID); err != nil {
		apiError(ctx, "User does not exist")
		return
	}

	if sec, err := model.RetrieveUserCredentials(user.ID); err != nil {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Failed to retrieve user credentials."),
		}, true)
		return
	} else if sec.Password != form.Password {
		util.ScrubbedPublicAPIJSON(ctx, gin.H{
			"error": fmt.Sprintf("Provided password is not the current user password."),
		}, true)
		return

	}

	if err := user.UpdatePassword(form.Password, form.NewPassword); err != nil {
		apiError(ctx, fmt.Sprintf("Failed to change password. %s", err))
		return
	}

	go cron.SendPasswordChangedNotification(user)
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"status": "changed",
	}, true)

}
