package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"iam-apiserver/internal/pkg/redis"
)

// Redis pub/sub events.
const (
	RedisPubSubChannel  = "iam.notifications"
	NoticePolicyChanged = "PolicyChanged"
	NoticeSecretChanged = "SecretChanged"
)

// Publish publishes a event to specified redis channel when policy/secret changed.
func Publish() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// ignores when policy/secret not changed
		if c.Writer.Status() != http.StatusOK || c.Request.Method == "GET" {
			return
		}

		var message string
		if message = getPublishMsg(c.Request.URL.Path); message == "" {
			return
		}

		if err := redis.Client().Publish(c, RedisPubSubChannel, message).Err(); err != nil {
			log.Printf("Publish error: %v", err)
		}
		log.Printf("Publish: %s - %s", RedisPubSubChannel, message)
	}
}

func getPublishMsg(URI string) string {
	var resource string
	pathSplit := strings.Split(URI, "/")
	if len(pathSplit) > 2 {
		resource = pathSplit[2]
	}

	var message string
	switch resource {
	case "policies":
		message = NoticePolicyChanged
	case "secrets":
		message = NoticeSecretChanged
	default:
	}

	return message
}
