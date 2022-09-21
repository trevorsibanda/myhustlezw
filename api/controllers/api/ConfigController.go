package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

type serviceConfig struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Title        string  `json:"title"`
	Price        float64 `json:"price"`
	Question     string  `json:"question"`
	ThankYou     string  `json:"thankyou"`
	Instructions string  `json:"instructions"`
	Quantity     int     `json:"quantity"`
}

var defaultServices = []serviceConfig{
	{
		Name:         "Custom",
		Description:  "Offer a service",
		Category:     "general",
		Title:        "",
		Price:        1.00,
		Question:     "Ask a question",
		ThankYou:     "Thank you for paying, I will contact you on the email you provided soon",
		Instructions: "Provide instructions here",
		Quantity:     10,
	},
}

//ServiceTemplates renders available templates for creating new services
func ServiceTemplates(ctx *gin.Context) {
	util.ScrubbedPublicAPIJSON(ctx, defaultServices, false)
}
