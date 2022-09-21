package cache

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

var (
	store                       *persistence.RedisStore
	cacheRedisHost              string
	cacheRedisPassword          string
	cacheRedisDefaultExpiration time.Duration = time.Second * 120
)

//InitCache initializes the cache store
func InitCache(r *gin.Engine) {
	cacheRedisHost = os.Getenv("CACHE_REDIS_HOST")
	cacheRedisPassword = os.Getenv("CACHE_REDIS_PASSWORD")
	store = persistence.NewRedisCache(cacheRedisHost, cacheRedisPassword, cacheRedisDefaultExpiration)

	if store == nil {
		panic("Failed to connect to redis caching server!")
	}
	store.Add("startup_time", time.Now(), time.Hour*480)
	var connectTime time.Time
	store.Get("startup_time", &connectTime)
	log.Println(fmt.Sprintf("Connected to redis cache store. Connected at %s", connectTime))
}

//Store returns the redis store
func Store() *persistence.RedisStore {
	return store
}
