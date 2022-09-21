package storage

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	model "github.com/trevorsibanda/myhustlezw/api/model"
)

var (
	myhustleScheme          = os.Getenv("MYHUSTLE_SCHEME")
	myhustleEndpoint        = os.Getenv("MYHUSTLE_DOMAIN")
	myhustleContentEndpoint = os.Getenv("MYHUSTLE_CONTENT_DOMAIN")
)

func fallbackImage(tpe string) string {
	return fmt.Sprintf("%s://%s/assets/img/placeholder%s.png", myhustleScheme, myhustleEndpoint, tpe)
}

//PreviewMediaURL generates a user image preview
func PreviewCreatorFileURL(media model.CreatorFile, width int, height int) (url string) {
	return
}

func GenerateTamperCode(fileId, fileName string) string {
	payload := fmt.Sprintf("%smy!huStle$zwnonce!!%s", strings.ToLower(fileId), strings.ToLower(fileName))
	hash := fmt.Sprintf("%x", (md5.Sum([]byte(payload))))
	return hash
}

func GetImageURL(sessionKey *string, file model.CreatorFile, width, height int, blurred bool) (uri string, err error) {

	opts := ".png"
	if width > 0 && height > 0 {
		blurImage := ""
		if blurred {
			blurImage = "_blurred"
		}
		opts = fmt.Sprintf("_%dx%d%s.png", width, height, blurImage)
	}

	switch file.Storage {
	case string(model.AssetFile):
		uri = fmt.Sprintf("%s://%s/assets/%s", myhustleScheme, myhustleEndpoint, file.Filename)
	case string(model.LocalFile):
		if blurred || sessionKey == nil {
			tamperProof := GenerateTamperCode(file.ID.Hex(), "file"+opts)
			uri = fmt.Sprintf("%s://%s/content/public/%s/%s/%s/file%s", myhustleScheme, myhustleContentEndpoint, file.Owner.Hex(), file.ID.Hex(), tamperProof, opts)
		} else {
			//	expires := time.Now().Add(time.Minute * 60 * 3).Unix()
			uri = fmt.Sprintf("%s://%s/content/session/%s/%s/%s/file%s", myhustleScheme, myhustleContentEndpoint, file.Owner.Hex(), file.ID.Hex(), file.AccessID(sessionKey), opts)
		}

	case string(model.YoutubeEmbed):
		uri = fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", file.Filename)
	case string(model.AWS):
		if parsed, err1 := url.Parse(awsS3Endpoint); err1 == nil {
			uri = fmt.Sprintf("https://%s.%s/%s", awsS3Bucket, parsed.Host, file.Filename)
		}
	default:
		uri = fallbackImage(file.Type)
	}
	return
}

//FilePreviewImageURL generates a preview url for a campaign file
func FilePreviewImageURL(sessionKey *string, file model.CreatorFile, notBlurred bool, width, height int) (url string) {
	var err error
	if url, err = GetImageURL(sessionKey, file, width, height, !notBlurred); err != nil {
		log.Printf("error generating url %v %v", file, err)
		url = fallbackImage("error")
	}
	return
}
