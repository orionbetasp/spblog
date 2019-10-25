package router

import (
	"fmt"
	"net/http"
	"spblog/conf"
	"spblog/controllers"
	"spblog/models"
	"spblog/util"
	"time"

	"go.uber.org/zap"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(controllers.SessionKey); uID != nil {
			user, err := models.GetUser(uID)
			if err == nil {
				c.Set(controllers.ContextUserKey, user)
			}
		}
		if conf.Con.GetBool("signup_enabled") {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AdminScopeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controllers.ContextUserKey); user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		util.Logger.Warn("User not authorized to visit " + c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!您可不是管理员哦！",
		})
		c.Abort()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controllers.ContextUserKey); user != nil {
			if _, ok := user.(*models.User); ok {
				c.Next()
				return
			}
		}
		util.Logger.Warn("User not authorized to visit " + c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!您可不是管理员哦！",
		})
		c.Abort()
	}
}

//校验时间和签名，取出参数
func signAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now().UnixNano() / 1000000

		c.Next()

		// 结束时间
		endTime := time.Now().UnixNano() / 1000000

		logMap := make(map[string]string)
		logMap["request_method"] = c.Request.Method
		logMap["request_uri"] = c.Request.RequestURI
		logMap["request_proto"] = c.Request.Proto
		logMap["request_ua"] = c.Request.UserAgent()
		logMap["request_referer"] = c.Request.Referer()
		logMap["request_client_ip"] = c.ClientIP()
		logMap["cost_time"] = fmt.Sprintf("%vms", endTime-startTime)
		util.Logger.Info("msg", zap.Any("info", logMap))
	}
}
