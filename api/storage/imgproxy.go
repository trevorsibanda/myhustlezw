package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
)

var (
	imgProxyKey      string = os.Getenv("IMGPROXY_KEY")
	imgProxySalt     string = os.Getenv("IMGPROXY_SALT")
	imgProxyEndpoint string = os.Getenv("IMGPROXY_ENDPOINT")
)

type ImageOptions struct {
	Resize    string //"fill"
	Gravity   string //"no"
	Enlarge   int    //:= 1
	Extension string //"png"
}

func ImgproxyURL(url string, width, height int, opts *ImageOptions) (path string, err error) {

	var keyBin, saltBin []byte

	if keyBin, err = hex.DecodeString(imgProxyKey); err != nil {
		err = fmt.Errorf("Key expected to be hex-encoded string")
		return
	}

	if saltBin, err = hex.DecodeString(imgProxySalt); err != nil {
		err = fmt.Errorf("Salt expected to be hex-encoded string")
		return
	}

	encodedURL := base64.RawURLEncoding.EncodeToString([]byte(url))

	path = fmt.Sprintf("/%s/%d/%d/%s/%d/%s.%s", opts.Resize, width, height, opts.Gravity, opts.Enlarge, encodedURL, opts.Extension)

	mac := hmac.New(sha256.New, keyBin)
	mac.Write(saltBin)
	mac.Write([]byte(path))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	path = fmt.Sprintf("%s%s%s", imgProxyEndpoint, signature, path)

	return
}
