package middleware

import (
	"net/http"
	"strings"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/gin-gonic/gin"

	"iam-apiserver/internal/pkg/redis"
)

// Redis pub/sub events.
const (
	channel            = "iam.notifications"
	eventPolicyChanged = "PolicyChanged"
	eventSecretChanged = "SecretChanged"
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

		if err := redis.Client().Publish(c, channel, message).Err(); err != nil {
			logger.L().X(c).Errorw("publish", "error", err.Error())
		}
		logger.L().X(c).Debugw("publish", "message", message, "method", c.Request.Method)
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
		message = eventPolicyChanged
	case "secrets":
		message = eventSecretChanged
	default:
	}

	return message
}
