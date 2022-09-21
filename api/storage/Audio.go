package storage

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	model "github.com/trevorsibanda/myhustlezw/api/model"
)

//ConvertAudioToAdaptiveHLS converts an audio to a streamable adaptive HLS
//This is done by downloading the video, running ffmpeg and then uploading the
//produced files
func ConvertAudioToAdaptiveHLS(audio model.CreatorFile) (hlsFile string, err error) {
	var inputFile, outputDir string

	out := fmt.Sprintf("%s/%s/", audio.Owner.Hex(), audio.ID.Hex())
	outputDir = LocalStoragePath(out)
	inputFile = fmt.Sprintf("%sfile%s", outputDir, audio.Extension)
	hlsFile = fmt.Sprintf("%sfile.m3u8", out)
	os.MkdirAll(outputDir, os.FileMode(0777))

	wd, _ := os.Getwd()
	command := exec.Command("bash", wd+"/scripts/audio_encoder/encoder.sh", inputFile, outputDir)
	command.Dir = wd + "/scripts/video_encoder/"
	log.Println(command.String())
	err = command.Run()
	return
}
