package controllers

import (
	"encoding/json"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/util"
)

type Page struct {
	Dict     map[string]string `json:"page" bson:"-" groups:"public"`
	Prefetch []string          `json:"prefetch" groups:"public"`
}

func (p *Page) MarshalJSON() ([]byte, error) {
	d := map[string]interface{}{
		"prefetch": strings.Join(p.Prefetch, "\n"),
	}

	for k, v := range p.Dict {
		nk := strings.ReplaceAll(k, "page_", "")
		d[nk] = v
	}

	return json.Marshal(d)
}

func (p *Page) SetTitle(title string) {
	p.Dict["page_title"] = title
}

func (p *Page) SetPageURL(url string) {
	p.Dict["page_url"] = url
}

func (p *Page) SetDescription(value string) {
	p.Dict["page_description"] = value
}

func (p *Page) SetImage(url string) {
	p.Dict["page_image"] = url
}

func (p *Page) SetImageAlt(alt string) {
	p.Dict["page_image_alt"] = alt
}

func (p *Page) SetPageAuthorURL(url string) {
	p.Dict["page_author_url"] = url
}

func (p *Page) SetTwitterUsername(username string) {
	p.Dict["twitter_username"] = username
}

func (p *Page) AddPrefetch(url string) {
	p.Prefetch = append(p.Prefetch, url)
}

func scrubbedPublicAPIJSON(ctx *gin.Context, data interface{}, loggedIn bool) {
	util.ScrubbedPublicAPIJSON(ctx, data, loggedIn)
}

func newPage(title string) (page Page) {
	page = Page{}
	page.Dict = map[string]string{}
	page.Prefetch = []string{}
	page.Dict["generated_at"] = time.Now().Format(time.UnixDate)
	page.SetTitle(title)
	return
}

//ReverseProxy transparently fowards a request
func ReverseProxy(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
