package model

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileType string

type FileRole string

const (
	//CoverImage is the cover image
	CoverImage FileRole = "cover"
	//ProfilePic is a users profile picture
	ProfilePic FileRole = "profile_spic"
	//ServicePreview is the preview file for a service
	ServicePreview FileRole = "service_preview"
	//Content can be photobook
	Content FileRole = "content"
)

var (
	previewURLGenerator func(*string, CreatorFile, bool, int, int) string
)

//RegisterPreviewURLGenerator registers a generator for preview urls
func RegisterPreviewURLGenerator(generator func(*string, CreatorFile, bool, int, int) string) {
	previewURLGenerator = generator
}

const (
	LocalFile       FileType = "local"
	AssetFile       FileType = "asset"
	AWS             FileType = "aws"
	RemoteFile      FileType = "remote"
	YoutubeEmbed    FileType = "youtube"
	SoundCloudEmbed FileType = "soundcloud"
)

//CreatorFile models files uploaded by a creator
type CreatorFile struct {
	ID               primitive.ObjectID `bson:"_id"  json:"_id" groups:"private"`
	Owner            primitive.ObjectID `bson:"owner_id"  json:"owner_id" groups:"private"`
	Deleted          bool               `bson:"deleted" json:"deleted" groups:"private"`
	OfferDownload    bool               `bson:"offer_download" json:"offer_download" groups:"public"`
	Created          time.Time          `bson:"created_at"  json:"created_at" groups:"private"`
	Filename         string             `bson:"filename"  json:"filename" groups:"protected"`
	Storage          string             `bson:"storage" json:"storage" groups:"private"`
	Extension        string             `bson:"ext" json:"ext" groups:"private"`
	ETag             string             `bson:"etag"  json:"etag" groups:"private"`
	Mime             string             `bson:"mime"  json:"mime" groups:"public"`
	OriginalFilename string             `bson:"original_filename"  json:"original_filename" groups:"private"`
	SizeBytes        uint               `bson:"size"  json:"size" groups:"public"`
	Campaign         primitive.ObjectID `bson:"campaign_id"  json:"campaign_id" groups:"private"`
	Type             string             `bson:"type"  json:"type" groups:"public"`
	Role             string             `bson:"role" json:"role" groups:"private"`
	Processed        bool               `bson:"processed" json:"processed" groups:"public"`
	Caption          string             `bson:"caption"  json:"caption" groups:"public"`
	URL              string             `bson:"-" json:"url" groups:"public"`
	Thumbnail        *string            `bson:"-" json:"thumbnail,omitempty" groups:"public"`
}

//SaveChanges saves the uploaded files and updates the changes
func (file CreatorFile) SaveChanges() (err error) {
	filter := bson.M{
		"_id": file.ID,
	}

	_, err = creatorFilesCollection().ReplaceOne(context.TODO(), filter, file)
	return
}

//Delete marks the file as deleted. The storage garbage collector will
//do the rest
func (file *CreatorFile) Delete() (err error) {
	file.Deleted = true
	return file.SaveChanges()
}

func (file *CreatorFile) AccessID(sessionKey *string) (aid string) {
	key := ""
	if sessionKey != nil {
		key = *sessionKey
	}
	payload := fmt.Sprintf("%snonce!_*%s%s", file.ID.Hex(), file.Owner.Hex(), key)
	aid = fmt.Sprintf("%x", md5.Sum([]byte(payload)))
	return
}

//SaveUploadToStorageChanges saves the uploaded files and updates the changes
func (file CreatorFile) SaveUploadToStorageChanges() (err error) {
	filter := bson.M{
		"_id": file.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"filename":  file.Filename,
			"storage":   file.Storage,
			"mime":      file.Mime,
			"processed": file.Processed,
			"etag":      file.ETag,
		},
	}

	_, err = creatorFilesCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//SetCampaign sets the campaign the file is associated with
func (file CreatorFile) SetCampaign(campaign *Campaign) (err error) {
	filter := bson.M{
		"_id": file.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"campaign_id": campaign.ID,
		},
	}

	_, err = creatorFilesCollection().UpdateOne(context.TODO(), filter, update)
	return
}

//GetFileByID gets a file given its id
func GetFileByID(id primitive.ObjectID, sessionKey *string, hasActiveSubscription bool) (file CreatorFile, err error) {
	filter := bson.M{
		"_id": id,
	}

	err = creatorFilesCollection().FindOne(context.TODO(), filter).Decode(&file)
	if err == nil {
		file.URL = file.PreviewURL(sessionKey, hasActiveSubscription, 480, 480)
		if file.Type == "image" || file.Type == "video" {
			thumb := file.PreviewURL(sessionKey, hasActiveSubscription, 120, 120)
			file.Thumbnail = &thumb
		}
	}
	return
}

//SaveUploadedFile saves an uploaded file and returns the document
func SaveUploadedFile(source *CreatorFile, pvtKey *string) (err error) {
	source.URL = source.PreviewURL(pvtKey, true, 0, 0)
	th := source.PreviewURL(pvtKey, true, 120, 120)
	source.Thumbnail = &th
	_, err = creatorFilesCollection().InsertOne(context.TODO(), source)

	return
}

//PreviewImageURL generates a preview image url
func (file CreatorFile) PreviewURL(sessionKey *string, hasAccess bool, width, height int) (url string) {
	url = previewURLGenerator(sessionKey, file, hasAccess, width, height)
	return
}
