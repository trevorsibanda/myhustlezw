package storage

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/trevorsibanda/myhustlezw/api/model"
)

//GenerateVideoPoster generates a high resolution poster from a video
//the generate video is uploaded and the url returned
func GenerateVideoPoster(video model.CreatorFile) (file string, err error) {
	//cmd := "-ss 3 -i \"%s\" -vf 'select=gt(scene\,0.4)' -frames:v 1 -vsync vfr -vf fps=fps=1/600 \"%s/%s.jpg\""
	var inputFile, outputFile string

	inputFile = LocalStoragePath(video.Filename)
	arr := strings.Split(video.Filename, ".")
	file = fmt.Sprintf("%s.png", arr[0])
	outputFile = LocalStoragePath(file)

	wd, _ := os.Getwd()
	log.Println(wd)
	command := exec.Command("bash", wd+"/scripts/video_poster/poster.sh", inputFile, outputFile)
	command.Dir = wd + "/scripts/video_poster/"
	err = command.Run()
	log.Println(command.String())
	return
}

func GetVideoStreamURL(video model.CreatorFile, hasSubscription bool) (uri string, err error) {

	//TODO: account for subscription n hls processing
	if video.Type != "video" && video.Type != "audio" {
		err = fmt.Errorf("Not a video or audio file")
		return
	}

	if !hasSubscription {
		uri = ""
		return
	}

	switch video.Storage {
	case string(model.AWS):
		if parsed, err1 := url.Parse(awsS3Endpoint); err1 == nil {
			uri = fmt.Sprintf("https://%s.%s/%s", awsS3Bucket, parsed.Host, video.Filename)
		}

	case string(model.LocalFile):
		expires := time.Now().Add(time.Minute * 600).Unix()
		payload := fmt.Sprintf("%s%d%s", "my!hustle!zw", expires, video.Campaign.Hex())
		expectedSignature := fmt.Sprintf("%x", md5.Sum([]byte(payload)))
		uri = fmt.Sprintf("%s://%s/stream/%s/%s/%s?e=%d&s=%s", myhustleScheme, myhustleContentEndpoint, "pvt", video.Campaign.Hex(), video.Filename, expires, expectedSignature)
	case string(model.YoutubeEmbed):
		uri = fmt.Sprintf("https://img.youtube.com/vi/%s/1.jpg", video.Filename)
	default:
		uri = "https://myhustle.co.zw/assets/img/blank.png"
	}
	return
}

//ConvertVideoToAdaptiveHLS converts a video to a streamable adaptive HLS
//This is done by downloading the video, running video2hls and then uploading the
//produced files
func ConvertVideoToAdaptiveHLS(video model.CreatorFile) (hlsFile string, err error) {
	var inputFile, outputDir string

	out := fmt.Sprintf("%s/%s/", video.Owner.Hex(), video.ID.Hex())
	outputDir = LocalStoragePath(out)
	inputFile = fmt.Sprintf("%sfile%s", outputDir, video.Extension)
	hlsFile = fmt.Sprintf("%sfile.m3u8", out)

	wd, _ := os.Getwd()
	command := exec.Command("bash", wd+"/scripts/video_encoder/runner.sh", inputFile, outputDir)
	command.Dir = wd + "/scripts/video_encoder/"
	log.Println(command.String())
	err = command.Run()

	//check if output file exists
	if err == nil {
		var file *os.File
		if file, err = os.Open(fmt.Sprintf("%s%s", outputDir, hlsFile)); err != nil {
			err = fmt.Errorf("Command ran successfully but file does not exist")
			command := exec.Command("bash", wd+"/scripts/video_encoder/fallback.sh", inputFile, outputDir)
			command.Dir = wd + "/scripts/video_encoder/"
			log.Println(command.String(), file)
			err = command.Run()
			return
		}
	}

	return
}
