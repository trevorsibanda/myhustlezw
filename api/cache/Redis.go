package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
)

var (
	redisClient   *redis.Client
	redisAddress  string
	redisHost     string
	redisUsername string
	redisPassword string
	redisDB       int
)

func NewRedisStore(name string) *redis.Client {
	redisAddress = os.Getenv("REDIS_STORE_ADDRESS")
	redisUsername = os.Getenv("REDIS_STORE_USERNAME")
	redisPassword = os.Getenv("REDIS_STORE_PASSWORD")
	redisDB, _ = strconv.Atoi(os.Getenv("REDIS_STORE_DB"))
	myhustleNodeID := os.Getenv("MYHUSTLE_NODE_ID")

	log.Printf("[%s] Connecting to redis store at %s with username %s and password %s", name, redisAddress, redisUsername, redisPassword)

	r := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Username: redisUsername,
		Password: redisPassword, // no password set
		DB:       redisDB,       // use default DB
	})
	now := time.Now()
	key := "startup_node_" + myhustleNodeID + "_" + name
	log.Printf("[%s] Redis store startup test at %s", name, now)
	r.Set(context.TODO(), key, now, time.Second*10)
	if timeResponse, err := r.Get(context.TODO(), key).Time(); err != nil {
		log.Fatalf("[%s] Redistore failed retrieve test with error %s", name, err)
		return nil
	} else {
		log.Printf("[%s] Redis wrote and retrieved %v", name, timeResponse)
	}

	return r
}

func InitRedisStore(r *gin.Engine) {
	redisClient = NewRedisStore("main")
	if redisClient == nil {
		panic("Failed to connect to redis store!")
	}
}

//Store returns the redis store
func RedisStore() *redis.Client {
	return redisClient
}
