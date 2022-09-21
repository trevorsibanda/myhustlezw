package main

import (
	"github.com/gin-gonic/gin"

	//"github.com/spf13/viper"
	"github.com/trevorsibanda/myhustlezw/api"
	"github.com/trevorsibanda/myhustlezw/api/cache"
	"github.com/trevorsibanda/myhustlezw/api/cron"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"github.com/trevorsibanda/myhustlezw/api/sessions"
	"github.com/trevorsibanda/myhustlezw/api/storage"
)

func main() {

	router := gin.Default()

	router.Use()

	model.InitDatabaseConnection()
	sessions.InitSecureWebsite(router)
	sessions.Init(router, "mylongasssessionstringhere")
	storage.InitStorageEngine()
	cache.InitCache(router)
	cache.InitRedisStore(router)
	cron.Init()

	api.Routes("v1", router)

	router.Run()
}
