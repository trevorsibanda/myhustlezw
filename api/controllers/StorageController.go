package controllers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/storage"
	"github.com/trevorsibanda/myhustlezw/api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var defaultImageOptions *storage.ImageOptions = &storage.ImageOptions{
	Gravity: "no",
	Enlarge: 1,
	Resize:  "fill",
}

//GenerateUploadURL an upload url for a particular file type
func GenerateUploadURL(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)

	mediaType, ok := ctx.Params.Get("type")
	if !ok {
		apiError(ctx, fmt.Sprintf("Invalid type"))
		return
	}

	var url string
	var headers map[string]string
	var err error

	var config storage.UploadConfig

	switch mediaType {
	case "video":
		config = storage.UploadConfig{
			ID:            primitive.NewObjectID(),
			AllowedMimes:  []string{"video/*"},
			DeleteFirst:   false,
			MaxUploadTime: time.Hour * 2,
		}
	case "audio":
		config = storage.UploadConfig{
			ID:            primitive.NewObjectID(),
			AllowedMimes:  []string{"audio/*"},
			DeleteFirst:   false,
			MaxUploadTime: time.Minute * 10,
		}
	case "image":
		config = storage.UploadConfig{
			ID:            primitive.NewObjectID(),
			AllowedMimes:  []string{"image/*"},
			DeleteFirst:   false,
			MaxUploadTime: time.Minute * 10,
		}
	case "other":
		config = storage.UploadConfig{
			ID:            primitive.NewObjectID(),
			AllowedMimes:  []string{"audio/*"},
			DeleteFirst:   false,
			MaxUploadTime: time.Hour * 1,
		}
	default:
		apiError(ctx, fmt.Sprintf("unknown type"))
		return
	}

	config.Owner = creator.ID
	config.Ext = "jpg"

	url, headers, err = storage.GenerateUploadToS3URL(config)
	if err != nil {
		log.Println("Failed to generated presigned upload url ", creator, " config: ", config, "err : ", err)
		apiError(ctx, fmt.Sprintf("Failed to generated presigned url"))
		return
	}
	util.ScrubbedPublicAPIJSON(ctx, gin.H{
		"url":     url,
		"method":  "PUT",
		"fields":  []string{},
		"headers": headers,
	}, false)

}

//UploadToStorage uploads the provided file to storage service
func UploadToStorage(ctx *gin.Context) {
	creator := ctx.Keys["creator"].(model.User)
	session, _ := sessions.GetVisitorSession(ctx)

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		apiError(ctx, fmt.Sprintf("invalid form"))
		return
	}

	//filename, err := storage.HandleUploadToS3(handle)
	var role string
	if roleP, ok := ctx.Params.Get("role"); ok {
		role = roleP
	}

	var validRoles = map[string]bool{
		"cover":           true,
		"profile_pic":     true,
		"content":         true,
		"service_preview": true,
		"feature":         true,
		"headline":        true,
	}
	if _, ok := validRoles[role]; !ok {
		apiError(ctx, fmt.Sprintf("invalid content role %s", role))
		return
	}

	var tpe = "image"
	//also process file here
	if t, ok := ctx.Params.Get("type"); ok {
		switch t {
		case "video":
			tpe = "video"
		case "audio":
			tpe = "audio"
		case "photobook":
			tpe = "image"
		case "image":
			tpe = "image"
		default:
			tpe = "other"
		}
	}

	var contentType string
	if file, e := fileHeader.Open(); e == nil {
		buffer := make([]byte, 512)
		file.Read(buffer)
		contentType = http.DetectContentType(buffer)
		defer file.Close()
	}

	if t, _ := ctx.Params.Get("type"); t == "image_video" {
		if strings.HasPrefix(contentType, "image/") {
			tpe = "image"
		} else {
			tpe = "video"
		}
	}

	//make sure ext is of allowed type
	ext := filepath.Ext(fileHeader.Filename)
	if tpe == "image" {
		ext = ".png"
	}

	file := model.CreatorFile{
		ID:               primitive.NewObjectID(),
		Owner:            creator.ID,
		Type:             tpe,
		Caption:          "",
		Processed:        false,
		Mime:             contentType,
		Created:          time.Now(),
		OriginalFilename: fileHeader.Filename,
		SizeBytes:        uint(fileHeader.Size),
		Storage:          string(model.LocalFile),
		Extension:        ext,
		Role:             role,
		Filename:         "",
	}

	dir := storage.LocalStoragePath(fmt.Sprintf("%s/%s", creator.ID.Hex(), file.ID.Hex()))
	if info, err := os.Lstat(dir); err != nil {
		if err := os.MkdirAll(dir, os.FileMode(0777)); err != nil {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to save the file. Failed to create directory",
			})
			return
		}
	} else if !info.IsDir() {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error. This user is not allowed to upload content.",
		})
		return
	}

	newFilename := fmt.Sprintf("%s/%s/file%s", creator.ID.Hex(), file.ID.Hex(), ext)

	if err := ctx.SaveUploadedFile(fileHeader, storage.LocalStoragePath(newFilename)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError,
			fmt.Errorf("Unable to save the file %v %v", storage.LocalStoragePath(newFilename), dir),
		)
		return
	}

	//save to db
	file.Filename = newFilename

	err = model.SaveUploadedFile(&file, &session.PrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save the file to database",
		})
		os.Remove(storage.LocalStoragePath(newFilename))
		return
	}

	//process file by role
	if file.Type == "image" {
		switch file.Role {
		case "profile_pic":
			creator.SetAvatar(file)
		case "cover":
			creator.SetCoverPicture(file)
		case "feature":
			creator.SetFeaturedContent(file)
		case "headline":
			creator.SetHeadlineSubscriptionImage(file)
		}
	}

	if file.Type == "video" {
		switch file.Role {
		case "feature":
			creator.SetFeaturedContent(file)
		}
	}

	//lets do it :)

	if file.Type == "video" {

		go func() {
			var videoPreview string
			videoPreview, err = storage.GenerateVideoPoster(file)
			if err != nil {
				fmt.Printf("failed to generate video preview with error %s \n", err)
			} else {

				//generate thumbnails
				image := file
				image.Type = "image"
				image.Mime = "image/png"
				image.Filename = videoPreview
				image.SizeBytes = 0
				image.Role = "video_preview"

				imageThumbnailGenerator(creator, image)

				//fmt.Printf("generated video preview at %s", videoPreview)
				//if storagePreview, err1 := storage.JustUploadFile(videoPreview, file.ID, videoPreview); err == nil {
				//	fmt.Printf("uploaded preview file to %s", storagePreview)
				//} else {
				//	fmt.Printf("failed to upload video preview to storage, %s", err1)
				//}
			}
		}()

	}

	util.ScrubbedPublicAPIJSON(ctx, file, true)
	//in the background, move the uploaded file to s3 and update db
	//only upload non content files

	if file.Type == "image" {
		go imageThumbnailGenerator(creator, file)
	}

	if file.Type == "video" {
		go processVideoEncode(creator, file)
	}

	if file.Type == "audio" {
		go processAudioEncode(creator, file)
	}

	//if file.Role != "content" {
	//	go uploadFileToRemoteServer(creator, file)
	//}

	return
}

func uploadFileToRemoteServer(creator model.User, file model.CreatorFile) {
	var err error
	newFilename, eTag, err := storage.HandleUploadToS3(creator, file)
	if err != nil {
		//todo delete aws file
		return
	}
	file.ETag = *eTag
	tmpFile := file.Filename
	file.Filename = newFilename
	file.Storage = string(model.AWS)
	file.Processed = true
	//save to db
	err = file.SaveUploadToStorageChanges()
	if err != nil {
		log.Println("Failed to save to database file we uploaded to s3", file, err)
	}
	err = os.Remove(storage.LocalStoragePath(tmpFile))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func processVideoEncode(creator model.User, file model.CreatorFile) {
	cron.ScheduleVideoEncoding(creator, file)
}

func processAudioEncode(creator model.User, file model.CreatorFile) {
	cron.ScheduleAudioEncoding(creator, file)
}

func imageThumbnailGenerator(creator model.User, file model.CreatorFile) {
	cron.ScheduleGenerateThumbnails(creator, file)
}

//AdaptiveImage returns a user uploaded image processed
func AdaptiveImage(ctx *gin.Context) {
	var err error
	creator := ctx.Keys["creator"].(model.User)

	var fileID primitive.ObjectID
	var width int
	var height int

	if id, ok := ctx.Params.Get("id"); ok {
		id = strings.TrimSuffix(id, ".png")
		fileID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			apiError(ctx, fmt.Sprintf("bad id"))
		}
	}

	if widthS, ok := ctx.Params.Get("width"); ok {
		width, err = strconv.Atoi(widthS)
		if err != nil {
			apiError(ctx, fmt.Sprintf("invalid width"))
		}
	}

	if heightS, ok := ctx.Params.Get("height"); ok {
		height, err = strconv.Atoi(heightS)
		if err != nil {
			apiError(ctx, fmt.Sprintf("invalid width"))
		}
	}

	opts := ""
	if width != 0 && height != 0 {
		opts = fmt.Sprintf("_%dx%d", width, height)
	}

	src := storage.LocalStoragePath(fmt.Sprintf("%s/%s/file%s.png", creator.ID.Hex(), fileID.Hex(), opts))
	file, err := os.Open(src)
	if err != nil {
		if opts != "" {
			src = storage.LocalStoragePath(fmt.Sprintf("%s/%s/file.png", creator.ID.Hex(), fileID.Hex()))
			if file, err = os.Open(src); err != nil {
				src = "./static/dashboard/public/img/placeholder.png"
				//apiError(ctx, fmt.Sprintf("Default file does not exist %v", err))
				//return
			}
		} else {
			apiError(ctx, fmt.Sprintf("File does not exist"))
			return
		}
	}

	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to get file info"))
		return
	}

	if stat.IsDir() {
		apiError(ctx, fmt.Sprintf("File is a directory"))
		return
	}

	extraHeaders := map[string]string{}
	ctx.DataFromReader(http.StatusOK, stat.Size(), "image/png", bufio.NewReader(file), extraHeaders)
}

//ServeContent serves the content to the user
// https://myhustle.co.zw/content/pub/1234567890123456790123456789012/09876543210987654321098765432109/file.png?e=764755765765&s=123456789096546465464545646
// https://myhustle.co.zw/content/pvt/1234567890123456790123456789012/09876543210987654321098765432109_file_480x480.png?e=764755765765&s=123456789096546465464545646&cc=211242143324234
func ServeContent(ctx *gin.Context) {
	session, _ := sessions.GetVisitorSession(ctx)

	permissions, _ := ctx.Params.Get("acl")

	if permissions != "public" && permissions != "session" {
		apiError(ctx, "Invalid permissions")
		return
	}

	user, _ := ctx.Params.Get("user")
	userId, _ := primitive.ObjectIDFromHex(user)
	if userId.IsZero() {
		apiError(ctx, "Invalid user")
		return
	}
	content, _ := ctx.Params.Get("contentid")
	contentId, _ := primitive.ObjectIDFromHex(content)
	if contentId.IsZero() {
		apiError(ctx, "no content specified")
		return
	}

	filename, _ := ctx.Params.Get("filename")

	//check if user can view content
	tamperCode := ctx.Param("tamper")
	lfile := model.CreatorFile{
		ID:    contentId,
		Owner: userId,
	}

	acl := "public"
	maxAge := 604800

	switch permissions {
	case "public":
		if tamperCode != storage.GenerateTamperCode(contentId.Hex(), filename) {
			log.Printf("%s %s %s %s", content, filename, tamperCode, storage.GenerateTamperCode(content, filename))
			apiError(ctx, "Tamper code check failed")
			return
		}

	case "session":
		if tamperCode != lfile.AccessID(&(session.PrivateKey)) {
			apiError(ctx, "You do not have permission to view this content")
			return
		}
		acl = "private"
	}

	ctx.Header("Cache-Control", fmt.Sprintf("%s, max-age=%d", acl, maxAge))
	ctx.Header("ETag", tamperCode)

	if match := ctx.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, tamperCode) {
			ctx.Status(http.StatusNotModified)
			return
		}
	}

	src := storage.LocalStoragePath(fmt.Sprintf("%s/%s/%s", user, content, filename))
	file, err := os.Open(src)
	if err != nil {
		if permissions != "public" {
			src = "./static/dashboard/public/assets/img/locked.jpeg"
			file, err = os.Open(src)
		} else if strings.HasSuffix(filename, ".png") {
			src = storage.LocalStoragePath(fmt.Sprintf("%s/%s/file.png", user, content))
			if file, err = os.Open(src); err != nil {
				src = "./static/dashboard/public/assets/img/placeholder.png"
				file, err = os.Open(src)
				if err != nil {
					apiError(ctx, fmt.Sprintf("Default does not exist %v", err))
					return
				}
				//apiError(ctx, fmt.Sprintf("Default file does not exist %v", err))
				//return
			}
		} else {
			apiError(ctx, fmt.Sprintf("File does not exist"))
			return
		}
	}

	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to get file info"))
		return
	}

	if stat.IsDir() {
		apiError(ctx, fmt.Sprintf("File is a directory"))
		return
	}

	extraHeaders := map[string]string{}
	ctx.DataFromReader(http.StatusOK, stat.Size(), "image/png", bufio.NewReader(file), extraHeaders)

}

//ServeDownload offers a file download
func ServeDownload(ctx *gin.Context) {

}

//StreamMedia is the endpoint for streaming video and audio files
func StreamMedia(ctx *gin.Context) {

	//permissions, _ := ctx.Params.Get("acl")

	user, _ := ctx.Params.Get("user")
	userId, _ := primitive.ObjectIDFromHex(user)
	if userId.IsZero() {
		apiError(ctx, "Invalid user")
		return
	}
	content, _ := ctx.Params.Get("content")
	contentId, _ := primitive.ObjectIDFromHex(content)
	if contentId.IsZero() {
		apiError(ctx, "no content specified")
		return
	}

	campaign, _ := ctx.Params.Get("content")
	campaignId, _ := primitive.ObjectIDFromHex(campaign)
	if campaignId.IsZero() {
		apiError(ctx, "no campaign specified")
		return
	}

	filename, _ := ctx.Params.Get("filename")

	//if time.Now().Unix() > int64(expires) {
	//	ctx.JSON(http.StatusForbidden, gin.H{
	//		"error": "This link has expired",
	//	})
	//}
	//log.Println(permissions, campaign, user, content, expires)
	//payload := fmt.Sprintf("%s%d%s", "my!hustle!zw", expires, campaign)
	//expectedSignature := fmt.Sprintf("%x", md5.Sum([]byte(payload)))
	//if expectedSignature != signature {
	//apiError(ctx, "signatures do not match. Tampering detected")
	//return
	//}
	log.Println("filename is ", filename)

	src := storage.LocalStoragePath(fmt.Sprintf("%s/%s/%s", user, content, filename))
	file, err := os.Open(src)
	if err != nil {
		if strings.HasSuffix(filename, ".png") {
			src = storage.LocalStoragePath(fmt.Sprintf("%s/%s/%s", user, content, filename))
			if file, err = os.Open(src); err != nil {
				src = "./static/dashboard/public/assets/videos/" + filename
				file, err = os.Open(src)
				if err != nil {
					apiError(ctx, fmt.Sprintf("Default does not exist %v", err))
					return
				}
				//apiError(ctx, fmt.Sprintf("Default file does not exist %v", err))
				//return
			}
		} else {
			apiError(ctx, fmt.Sprintf("File does not exist"))
			return
		}
	}

	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil {
		apiError(ctx, fmt.Sprintf("Failed to get file info"))
		return
	}

	if stat.IsDir() {
		apiError(ctx, fmt.Sprintf("File is a directory"))
		return
	}

	mime, _ := util.GetFileContentType(file)
	file.Seek(0, 0)

	extraHeaders := map[string]string{}
	ctx.DataFromReader(http.StatusOK, stat.Size(), mime, bufio.NewReader(file), extraHeaders)

}
