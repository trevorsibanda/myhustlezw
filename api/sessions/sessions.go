package sessions

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionInfo struct {
	IP       string
	Browser  string
	Location string
}

func (i SessionInfo) Dictionary() map[string]string {
	return map[string]string{
		"IP Address": i.IP,
		"Browser":    i.Browser,
		"Location":   i.Location,
	}
}

//VisitorSession models a visitor's session
type VisitorSession struct {
	Token       string               `json:"token" groups:"private"`
	HasCreator  bool                 `json:"has_creator" groups:"private"`
	LastActive  time.Time            `json:"last_active" groups:"private"`
	User        model.User           `json:"user" groups:"private"`
	PrivateKey  string               `json:"key" groups:"private"`
	PaidContent []primitive.ObjectID `json:"paid_content" groups:"public"`
	Info        SessionInfo          `json:"-" groups:"protected"`
}

//NewVisitorSession creates a new session for a visitor
func NewVisitorSession(ctx *gin.Context) (session VisitorSession) {
	sessions.Default(ctx).Clear()

	session = VisitorSession{
		Token:       generateAuthToken(primitive.NewObjectID()),
		LastActive:  time.Now(),
		PaidContent: []primitive.ObjectID{},
	}
	session.PrivateKey = fmt.Sprintf("%x", md5.Sum([]byte(session.Token)))
	session.UpdateInfo(ctx)
	//write

	sess := sessions.Default(ctx)

	parsed, _ := json.Marshal(session)
	sess.Set("visitor_session", string(parsed))
	sess.Save()
	return
}

//LastActiveNow sets the visitor's session to last seen now
func (vs *VisitorSession) LastActiveNow() {
	vs.LastActive = time.Now()
}

//GinSession returns the underlying gin session
func (vs *VisitorSession) GinSession(ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}

//Save saves the current state of the visitorsession
func (vs VisitorSession) Save(ctx *gin.Context) error {
	sess := vs.GinSession(ctx)
	var err error
	var parsed []byte
	if parsed, err = json.Marshal(vs); err != nil {
		log.Printf("Failed to save visitor session %v", err)
	} else {
		//Save saves the current state of the visitorsession
		sess.Set("visitor_session", string(parsed))
		err = sess.Save()
	}
	return err
}

func (vsession *VisitorSession) UpdateInfo(ctx *gin.Context) {
	location := "Unknown"
	vsession.Info = SessionInfo{
		IP:       ctx.ClientIP(),
		Browser:  ctx.Request.UserAgent(),
		Location: location,
	}
}

//UpdateGrantContentAccess updates the session's paid content
func (vsession *VisitorSession) UpdateGrantContentAccess(ctx *gin.Context, id *primitive.ObjectID) (err error) {
	var ids []primitive.ObjectID
	if vsession.User.ID.IsZero() && id != nil {
		ids = append(vsession.PaidContent, *id)
	} else if !vsession.User.ID.IsZero() && id == nil {
		ids, _ = vsession.User.ListPaidCampaigns()
	}

	if id != nil {
		ids = append(ids, *id)
	}
	vsession.PaidContent = ids
	if err = vsession.Save(ctx); err != nil {
		log.Println(fmt.Sprintf("Failed to save visitor session %v", err))
		return
	}
	return
}

//Init initializes the http session session store
func Init(r *gin.Engine, sessionSecret string) {
	secret := os.Getenv("SESSIONS_REDIS_SECRET")
	if store, err := redis.NewStore(20, "tcp", os.Getenv("SESSIONS_REDIS_ADDRESS"), os.Getenv("SESSIONS_REDIS_PASSWORD"), []byte(secret)); err != nil {
		panic(fmt.Sprintf("Failed to create new redis session store %s ", err))
	} else {
		store.Options(sessions.Options{
			Path:   "/",
			MaxAge: int(time.Hour * 6),
		})
		r.Use(sessions.Sessions("session", store))
	}
	log.Println("Connected to redis sessions store")
}

//GetVisitorSession returns the current session's data from the auth_token
func GetVisitorSession(ctx *gin.Context) (session VisitorSession, err error) {
	sess := sessions.Default(ctx)
	raw := sess.Get("visitor_session")

	if raw != nil {
		unparsed := raw.(string)
		err = json.Unmarshal([]byte(unparsed), &session)
		session.UpdateInfo(ctx)
	} else {
		err = fmt.Errorf("Failed to get visitor session")
		return
	}
	return
}

//generateAuthToken generates an auth token given a user id
func generateAuthToken(userID primitive.ObjectID) (token string) {
	token = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%d", userID.Hex(), time.Now().String(), rand.Int31()))))
	return
}
