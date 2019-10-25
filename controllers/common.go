package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"spblog/conf"
	"spblog/models"
	"spblog/util"
	"strings"

	"github.com/denisbakhtin/sitemap"
	"github.com/gin-gonic/gin"
)

const (
	SessionKey         = "UserID"       // session key
	ContextUserKey     = "User"         // context user key
	SessionGithubState = "GITHUB_STATE" // github state session key
	SessionCaptcha     = "GIN_CAPTCHA"  // captcha session key
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
}

func sendMail(to, subject, body string) error {
	return util.SendToMail(conf.Con.GetString("smtp.smtp_username"),
		conf.Con.GetString("smtp.smtp_password"), conf.Con.GetString("smtp.smtp_host"), to, subject, body, "html")
}

func NotifyEmail(subject, body string) error {
	notifyEmailsStr := conf.Con.GetString("notify_emails")
	if notifyEmailsStr != "" {
		notifyEmails := strings.Split(notifyEmailsStr, ";")
		emails := make([]string, 0)
		for _, email := range notifyEmails {
			if email != "" {
				emails = append(emails, email)
			}
		}
		if len(emails) > 0 {
			return sendMail(strings.Join(emails, ";"), subject, body)
		}
	}
	return nil
}

func CreateXMLSitemap() {
	folder := path.Join(conf.Con.GetString("public"), "sitemap")
	_ = os.MkdirAll(folder, os.ModePerm)
	domain := conf.Con.GetString("domain")
	now := util.GetCurrentTime()
	items := make([]sitemap.Item, 0)

	items = append(items, sitemap.Item{
		Loc:        domain,
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	posts, err := models.ListPublishedPost("", 0, 0)
	if err == nil {
		for _, post := range posts {
			items = append(items, sitemap.Item{
				Loc:        fmt.Sprintf("%s/post/%d", domain, post.ID),
				LastMod:    post.UpdatedAt,
				Changefreq: "weekly",
				Priority:   0.9,
			})
		}
	}

	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/static/sitemap/"); err != nil {
		return
	}
}

func writeJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}
