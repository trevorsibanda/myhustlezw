package payments

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/trevorsibanda/myhustlezw/api/model"
)

var (
	paymentAPIEndpoint         = os.Getenv("PAYMENT_API_ENDPOINT")
	paymentApiCallbackEndpoint = os.Getenv("PAYMENT_API_CALLBACK_ENDPOINT")
	paymentAPISecret           = os.Getenv("PAYMENT_API_SECRET")
	client                     http.Client
)

type EcocashPollResponse struct {
	Reference       string  `json:"reference"`
	PaynowReference string  `json:"paynowreference"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status" `
	PollURL         string  `json:"pollurl"`
	Hash            string  `json:"hash"`
}

func (r *EcocashPollResponse) FromFormData(data string) error {
	values, err := url.ParseQuery(data)
	r.Status = values.Get("status")
	r.Hash = values.Get("hash")
	r.PaynowReference = values.Get("paynowreference")
	r.Reference = values.Get("reference")
	r.Amount, _ = strconv.ParseFloat(values.Get("amount"), 32)
	r.PollURL = values.Get("pollurl")
	return err
}

func PaymentURI(action string, payment model.PendingPayment) string {
	return fmt.Sprintf("%s%s/%s", paymentAPIEndpoint, action, payment.Gateway)
}

func PaymentAPIChallenge(id string) (sha1_hash string) {
	h := sha512.New()
	h.Write([]byte(id + paymentAPISecret + id))
	sha1_hash = hex.EncodeToString(h.Sum(nil))
	return
}

func ApiRequest(endpoint string, payment model.PendingPayment) (result []byte, err error) {
	payment.RedirectURL = paymentApiCallbackEndpoint + "/" + string(payment.Gateway)
	data, _ := json.Marshal(payment)
	response, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))

	if err == nil {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err == nil {
			result = body
		}
	}

	return
}
