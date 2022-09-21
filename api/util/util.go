package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/dustin/go-humanize/english"
	"github.com/gin-gonic/gin"
	"github.com/liip/sheriff"
)

var (
	publicDataSheriffOption = &sheriff.Options{
		Groups: []string{"public"},
	}

	authenticatedDataSheriffOption = &sheriff.Options{
		Groups: []string{"public", "private", "authenticated"},
	}

	noPregenCache = (os.Getenv("MYHUSTLE_NO_PREGEN_CACHE") == "1")
)

func ApiError(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusOK, gin.H{
		"error": reason,
	})
}

func ScrubbedPublicAPIJSON(ctx *gin.Context, data interface{}, loggedIn bool) {
	opts := publicDataSheriffOption
	if loggedIn {
		opts = authenticatedDataSheriffOption
	}
	data, err := sheriff.Marshal(opts, data)
	if err != nil {
		ApiError(ctx, "Sheriff failed to sanitize data to JSON.")
		log.Printf("[sheriff] Failed to sanitize data to JSON. Error: %v Data: %v", err, data)
		return
	}

	ctx.JSON(200, data)
}

func GenerateMeta(dict map[string]string) string {
	generated := metaTemplate
	for k, v := range dict {
		generated = strings.ReplaceAll(generated, fmt.Sprintf("{{%s}}", k), escape(v))
	}
	return generated
}

func GeneratePageWithData(prefix, html string, data interface{}) string {
	tag := "<pregenerated></pregenerated>"
	if noPregenCache {
		tag = "<nope><nope/>"
	}
	json, _ := json.Marshal(ScrubPublic(data))
	return strings.ReplaceAll(html, tag, fmt.Sprintf(pregenTemplate, prefix, json, prefix))
}

func escape(s string) string {
	return strings.ReplaceAll(s, "\"", `\"`)
}

func ScrubPublic(data interface{}) (result interface{}) {
	var err error
	opts := publicDataSheriffOption
	result, err = sheriff.Marshal(opts, data)
	if err != nil {
		result = nil
	}
	return
}

func ToDict(values ...interface{}) (map[string]interface{}, error) {
	if len(values) == 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]interface{})

	for i := 0; i < len(values); i++ {
		key, isset := values[i].(string)
		if !isset {
			if reflect.TypeOf(values[i]).Kind() == reflect.Map {
				m := values[i].(map[string]interface{})
				for i, v := range m {
					dict[i] = v
				}
			} else {
				return nil, errors.New("dict values must be maps")
			}
		} else {
			i++
			if i == len(values) {
				return nil, errors.New("specify the key for non array values")
			}
			dict[key] = values[i]
		}

	}
	return dict, nil
}

func ReadableBytes(size uint) string {
	return humanize.Bytes(uint64(size))
}

func PluralOf(n int, item string) string {
	return english.Plural(n, item, "")
}

func OrdinalOf(n int) string {
	return humanize.Ordinal(n)
}

func TimeAgo(t time.Time) string {
	return humanize.Time(t)
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
