package cron

import (
	"fmt"

	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
)

//OnSignupVerifiedPhone is logic to be executed once a phone number was verified from signup
func OnSignupVerifiedPhone(user model.User, session sessions.VisitorSession) (err error) {

	err = model.UpdatePhoneNumberAndPaymentsToVerified(user)

	SendSignupEmail(user, session)
	SendGettingStartedEmail(user)
	return

}

func OnEmailVerified(user model.User) (err error) {
	OnEmailChanged(user)
	return
}

func OnEmailChanged(user model.User) (err error) {
	err = model.UpdateEmailToVerified(user)
	SendEmailChangedEmail(user.Email, user)
	return
}

func OnPhoneChanged(user model.User) (err error) {
	err = model.UpdatePhoneToVerified(user)
	SendPhoneNumberChangedEmail(user)
	return
}

func OnResetPasswordPhone(user model.User) (err error) {
	return
}

//OnAddSupportMessageOnPaid is evoked when a support message is left
func OnAddSupportMessageOnPaid(supporter model.User, creator model.User, payment model.PendingPayment, support model.CreatorSupport) (err error) {
	//notify supporter of success paid
	buyItem := fmt.Sprintf("one %s", creator.Page.DonationItemName)
	if payment.Items > 1 {
		buyItem = fmt.Sprintf("%d %ss", payment.Items, creator.Page.DonationItemName)
	}
	var amount string
	switch payment.Currency {
	case model.USD:
		amount = FormatAsUSD(payment.Price)
	case model.ZWL:
		amount = FormatAsZWL(payment.Price)
	}
	SMSmsg := fmt.Sprintf("You have successfully bought %s %s for %s", creator.Username, buyItem, amount)
	emailMsg := fmt.Sprintf("Hey %s,\n%s", payment.Fullname, SMSmsg)
	notifyPayerSMS := model.SMSNotification{
		PhoneNumber: support.PhoneNumber,
	}
	notifyPayerSMS.Message = SMSmsg
	notifyPayerSMS.Owner = creator.ID

	notifyPayerEmail := model.EmailNotification{
		Title:      fmt.Sprintf("You have bought @%s %s", creator.Username, buyItem),
		Template:   "support_message",
		ActionURL:  creator.URL(),
		ActionName: fmt.Sprintf("View %s's page", creator.Username),
		Dictionary: map[string]string{
			"Amount Paid": Format(string(payment.Currency), payment.Price),
			"Gateway":     string(payment.Gateway),
			"Payment ID":  payment.ID.Hex(),
		},
	}
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = creator.ID

	_ = SendSMSNotification(&notifyPayerSMS)

	DispatchEmail(supporter, notifyPayerEmail, false)

	if creator.Notifications.WalletCredit {
		emailMsg = fmt.Sprintf("Congrats!, \n%s bought you %s. You can view this payment in detail in your dashboard", support.DisplayName, buyItem)
		notifyPayeeEmail := model.EmailNotification{
			Title:      fmt.Sprintf("%s bought you %s", support.DisplayName, buyItem),
			Template:   "support_message_creator",
			ActionURL:  creator.URL(),
			ActionName: "Go to your dashboard",
			Dictionary: map[string]string{
				"Amount Paid": Format(string(payment.Currency), payment.Price),
				"Gateway":     string(payment.Gateway),
				"Payment ID":  payment.ID.Hex(),
			},
		}
		notifyPayeeEmail.Email = creator.Email
		notifyPayeeEmail.Message = emailMsg
		notifyPayeeEmail.Owner = creator.ID
		DispatchEmail(creator, notifyPayeeEmail, false)
	}
	return
}

//OnAddServiceRequest is evoked when a service request is granted
func OnAddServiceRequest(supporter model.User, creator model.User, service model.Campaign, payment model.PendingPayment, support model.CreatorSupport) (err error) {
	//notify supporter of success paid
	var amount string
	switch payment.Currency {
	case model.USD:
		amount = FormatAsUSD(payment.Price)
	case model.ZWL:
		amount = FormatAsZWL(payment.Price)
	}
	SMSmsg := fmt.Sprintf("You have successfully paid %s for \"%s\" offered by @%s. We sent you an email with more details.", amount, service.Title, creator.Username)
	emailMsg := fmt.Sprintf("Hey %s,\n%s", payment.Fullname, SMSmsg)
	notifyPayerSMS := model.SMSNotification{
		PhoneNumber: support.PhoneNumber,
	}
	notifyPayerSMS.Message = SMSmsg
	notifyPayerSMS.Owner = creator.ID

	dict := support.Form.Dict()
	dict["Amount Paid"] = Format(string(payment.Currency), payment.Price)
	dict["Gateway"] = string(payment.Gateway)
	dict["Payment ID"] = payment.ID.Hex()

	//TODO: send access link
	notifyPayerEmail := model.EmailNotification{
		Title:      fmt.Sprintf("You have paid %s for a service offered by @%s ", amount, creator.Username),
		Template:   "service_request",
		ActionURL:  service.URL(creator.Username, &support.ID, &support.UnlockCode),
		ActionName: fmt.Sprintf("View %s's page", creator.Username),
		Dictionary: dict,
	}
	emailMsg += "\nIf your order is not marked as fulfilled in 72 hours we will refund it within 72hours."
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = creator.ID

	if err = SendSMSNotification(&notifyPayerSMS); err != nil {
		return
	}

	DispatchEmail(supporter, notifyPayerEmail, false)

	emailMsg = fmt.Sprintf("Congrats!, \n%s paid for your \"%s\" service offering . You can view this payment in detail in your dashboard", support.DisplayName, service.Title)
	emailMsg += "\nRemember you have up to 72 hours to go to your dashboard to honor this request or the money will be refunded to the payer."
	notifyPayeeEmail := model.EmailNotification{
		Title:      fmt.Sprintf("%s paid for a service you are offering", support.DisplayName),
		Template:   "service_request_creator",
		ActionURL:  support.DashboardURL(),
		ActionName: "View order",
		Dictionary: dict,
	}
	notifyPayeeEmail.Email = creator.Email
	notifyPayeeEmail.Message = emailMsg
	notifyPayeeEmail.Owner = creator.ID
	DispatchEmail(creator, notifyPayeeEmail, false)

	if service.Form.QuantityAvailable == 0 {
		emailMsg = fmt.Sprintf("You will no longer be recieving any more orders for \"%s\". To accept more orders, edit the service to have more items available. Failure to do this will result in the service getting archived in seven days.", service.Title)
		notifyPayeeEmail := model.EmailNotification{
			Title:      fmt.Sprintf("No longer accepting orders for %s", service.Title),
			Template:   "notification",
			ActionURL:  support.DashboardURL(),
			ActionName: "Go to my account",
		}
		notifyPayeeEmail.Email = creator.Email
		notifyPayeeEmail.Message = emailMsg
		notifyPayeeEmail.Owner = creator.ID
		DispatchEmail(creator, notifyPayeeEmail, false)
	}
	return
}

//OnServiceFulfilled is evoked when a service has been fulfilled
func OnServiceFulfilled(creator model.User, support model.CreatorSupport, session sessions.VisitorSession) {
	SMSmsg := fmt.Sprintf("@%s has marked your order %s for \"%s\" as done. Check your email for more details.", creator.Username, support.ID.Hex(), support.ItemName)
	emailMsg := fmt.Sprintf("Hey %s,\n%s", support.DisplayName, SMSmsg)
	notifyPayerSMS := model.SMSNotification{
		PhoneNumber: support.PhoneNumber,
	}
	notifyPayerSMS.Message = SMSmsg
	notifyPayerSMS.Owner = creator.ID

	dict := support.Form.Dict()
	dict["Amount Paid"] = Format(support.Currency, support.Amount)
	dict["Order ID"] = support.ID.Hex()

	//TODO: send access link
	notifyPayerEmail := model.EmailNotification{
		Title:      fmt.Sprintf("@%s has fulfilled your order", creator.Username),
		Template:   "notification",
		ActionURL:  creator.URL(),
		ActionName: fmt.Sprintf("View @%s's page", creator.Username),
		Dictionary: dict,
	}
	emailMsg += "\nIf you have a dispute, please contact support or open a payment dispute."
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = creator.ID

	if err := SendSMSNotification(&notifyPayerSMS); err != nil {
		return
	}
	var supporter model.User
	if support.Supporter.IsZero() {
		supporter = model.AnonymousUser(support.Email, support.PhoneNumber, support.Form.Fullname)
	} else {
		supporter, _ = model.RetrieveCreatorByID(support.Supporter)
	}

	DispatchEmail(supporter, notifyPayerEmail, false)

	emailMsg = "The status of the order has been changed to fulfilled. The funds will be automatically released from escrow and you will be notified when you can withdraw."
	notifyPayeeEmail := model.EmailNotification{
		Title:      fmt.Sprintf("You fulfilled order %s by %s", support.ID.Hex(), support.DisplayName),
		Template:   "notification",
		ActionURL:  support.DashboardURL(),
		ActionName: "View order",
		Dictionary: dict,
	}
	notifyPayeeEmail.Email = creator.Email
	notifyPayeeEmail.Message = emailMsg
	notifyPayeeEmail.Owner = creator.ID
	DispatchEmail(creator, notifyPayeeEmail, false)
}

//OnServiceRefunded is evoked when a service has been fulfilled
func OnServiceRefunded(creator model.User, support model.CreatorSupport, session sessions.VisitorSession) {
	SMSmsg := fmt.Sprintf("@%s has refunded your order for \"%s\". The refund will be processed with 72 hours. Check your email for more details.", creator.Username, support.ItemName)
	emailMsg := fmt.Sprintf("Hey %s,\n%s", support.DisplayName, SMSmsg)
	notifyPayerSMS := model.SMSNotification{
		PhoneNumber: support.PhoneNumber,
	}
	notifyPayerSMS.Message = SMSmsg
	notifyPayerSMS.Owner = creator.ID

	dict := support.Form.Dict()
	dict["Amount Paid"] = Format(support.Currency, support.Amount)
	dict["Order ID"] = support.ID.Hex()

	//TODO: send access link
	notifyPayerEmail := model.EmailNotification{
		Title:      fmt.Sprintf("@%s has refunded your order", creator.Username),
		Template:   "notification",
		ActionURL:  creator.URL(),
		ActionName: fmt.Sprintf("View @%s's page", creator.Username),
		Dictionary: dict,
	}
	emailMsg += "\nIf you have a dispute, please contact support or open a payment dispute."
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = creator.ID

	if err := SendSMSNotification(&notifyPayerSMS); err != nil {
		return
	}
	var supporter model.User
	if support.Supporter.IsZero() {
		supporter = model.AnonymousUser(support.Email, support.PhoneNumber, support.Form.Fullname)
	} else {
		supporter, _ = model.RetrieveCreatorByID(support.Supporter)
	}

	DispatchEmail(supporter, notifyPayerEmail, false)

	emailMsg = "Your request to refund the order has been received. We will handle everything from here and make sure the money gets paid back into the payer's account."
	notifyPayeeEmail := model.EmailNotification{
		Title:      fmt.Sprintf("You refunded order %s by %s", support.ID.Hex(), support.DisplayName),
		Template:   "notification",
		ActionURL:  support.DashboardURL(),
		ActionName: "View order",
		Dictionary: dict,
	}
	notifyPayeeEmail.Email = creator.Email
	notifyPayeeEmail.Message = emailMsg
	notifyPayeeEmail.Owner = creator.ID
	DispatchEmail(creator, notifyPayeeEmail, false)
}

//OnPayPerViewAccessGranted is evoked when a pay per view access is granted
func OnPayPerViewAccessGranted(supporter model.User, campaign model.Campaign, creator model.User, payment model.PendingPayment, support model.CreatorSupport) (err error) {
	//notify supporter of success paid
	var amount string
	switch payment.Currency {
	case model.USD:
		amount = FormatAsUSD(payment.Price)
	case model.ZWL:
		amount = FormatAsZWL(payment.Price)
	}
	SMSmsg := fmt.Sprintf("%s payment to view content by @%s successful. A link for future access was sent to your email. ", amount, creator.Username)
	emailMsg := fmt.Sprintf("Hey %s,\n%s", payment.Fullname, SMSmsg)
	notifyPayerSMS := model.SMSNotification{
		PhoneNumber: support.PhoneNumber,
	}
	notifyPayerSMS.Message = SMSmsg
	notifyPayerSMS.Owner = creator.ID

	notifyPayerEmail := model.EmailNotification{
		Title:      fmt.Sprintf("You have paid %s to access @%s's content", amount, creator.Username),
		Template:   "pay_per_view",
		ActionURL:  campaign.URL(creator.Username, &support.ID, &support.UnlockCode),
		ActionName: "View content",
		Dictionary: map[string]string{
			"Amount Paid": Format(string(payment.Currency), payment.Price),
			"Gateway":     string(payment.Gateway),
			"Payment ID":  payment.ID.Hex(),
		},
	}
	notifyPayerEmail.Message = emailMsg
	notifyPayerEmail.Owner = creator.ID

	_ = SendSMSNotification(&notifyPayerSMS)

	DispatchEmail(supporter, notifyPayerEmail, false)

	if creator.Notifications.WalletCredit {
		emailMsg = fmt.Sprintf("Congrats!, \n%s paid to access your content - %s. You can view this payment in detail in your dashboard", support.DisplayName, campaign.Title)
		notifyPayeeEmail := model.EmailNotification{
			Title:      fmt.Sprintf("%s paid to view %s - %s", support.DisplayName, campaign.Type, campaign.Title),
			Template:   "support_message_creator",
			ActionURL:  creator.URL(),
			ActionName: "Go to your dashboard",
			Dictionary: map[string]string{
				"Amount Paid": Format(string(payment.Currency), payment.Price),
				"Gateway":     string(payment.Gateway),
				"Payment ID":  payment.ID.Hex(),
			},
		}
		notifyPayeeEmail.Email = creator.Email
		notifyPayeeEmail.Message = emailMsg
		notifyPayeeEmail.Owner = creator.ID
		DispatchEmail(creator, notifyPayeeEmail, false)
	}
	return
}
