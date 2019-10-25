package controllers

import (
	"fmt"
	"net/http"
	"spblog/conf"
	"spblog/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthGet(c *gin.Context) {
	authType := c.Param("authType")

	session := sessions.Default(c)
	uuid := util.UUID()
	session.Delete(SessionGithubState)
	session.Set(SessionGithubState, uuid)
	_ = session.Save()

	authurl := "/signin"
	switch authType {
	case "github":
		authurl = fmt.Sprintf(conf.Con.GetString("github.github_authurl"), conf.Con.GetString("github.github_clientid"), uuid)
	case "weibo":
	case "qq":
	case "wechat":
	case "oschina":
	default:
	}
	util.Logger.Info(authurl)
	c.Redirect(http.StatusFound, authurl)
}
