package middleware

// import (
// 	"context"
// 	"encoding/json"
// 	"iam-apiserver/internal/pkg/cache"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// // Redis pub/sub events.
// const (
// 	RedisPubSubChannel  = "iam.notifications"
// 	NoticePolicyChanged = "PolicyChanged"
// 	NoticeSecretChanged = "SecretChanged"
// )

// // Publish publish a event to specified redis channel when policy/secret changed.
// func Publish() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next()

// 		if c.Writer.Status() != http.StatusOK {
// 			return
// 		}

// 		method := c.Request.Method

// 		var resource string

// 		pathSplit := strings.Split(c.Request.URL.Path, "/")
// 		if len(pathSplit) > 2 {
// 			resource = pathSplit[2]
// 		}

// 		switch resource {
// 		case "policies":
// 			notify(c, method, load.NoticePolicyChanged)
// 		case "secrets":
// 			notify(c, method, load.NoticeSecretChanged)
// 		default:
// 		}
// 	}
// }

// func notify(ctx context.Context, method string, command load.NotificationCommand) {
// 	switch method {
// 	case "POST", "PUT", "DELETE":
// 		redisStore := cache.Cache()
// 		message, _ := json.Marshal(load.Notification{Command: command})

// 		if err := redisStore.Publish(load.RedisPubSubChannel, string(message)); err != nil {
// 			log.Printf("publish redis message failed", "error", err.Error())
// 		}
// 		log.Printf("publish redis message", "method", method, "command", command)
// 	default:
// 	}
// }
