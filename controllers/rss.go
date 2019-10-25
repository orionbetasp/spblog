package controllers

import (
	"fmt"
	"spblog/conf"
	"spblog/models"
	"spblog/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

func RssGet(c *gin.Context) {
	now := util.GetCurrentTime()
	domain := conf.Con.GetString("domain")
	feed := &feeds.Feed{
		Title:       "spblog",
		Link:        &feeds.Link{Href: domain},
		Description: "spblog,talk about golang,java and so on.",
		Author:      &feeds.Author{Name: "Wangsongyan", Email: "wangsongyanlove@163.com"},
		Created:     now,
	}

	feed.Items = make([]*feeds.Item, 0)
	posts, err := models.ListPublishedPost("", 0, 0)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}

	for _, post := range posts {
		item := &feeds.Item{
			Id:          fmt.Sprintf("%s/post/%d", domain, post.ID),
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%d", domain, post.ID)},
			Description: string(post.Excerpt()),
			Created:     now,
		}
		feed.Items = append(feed.Items, item)
	}
	rss, err := feed.ToRss()
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	_, _ = c.Writer.WriteString(rss)
}
