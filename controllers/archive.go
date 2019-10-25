package controllers

import (
	"math"
	"net/http"
	"spblog/conf"
	"spblog/models"
	"spblog/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func ArchiveGet(c *gin.Context) {
	var (
		year      string
		month     string
		page      string
		pageIndex int
		pageSize  = conf.Con.GetInt("page_size")
		total     int
		err       error
		posts     []*models.Post
		policy    *bluemonday.Policy
	)
	year = c.Param("year")
	month = c.Param("month")
	page = c.Query("tag")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	posts, err = models.ListPostByArchive(year, month, pageIndex, pageSize)
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByArchive(year, month)
	if err != nil {
		util.Logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(blackfriday.Run([]byte(post.Body))))
	}
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts":           posts,
		"tags":            models.MustListTag(),
		"archives":        models.MustListPostArchives(),
		"links":           models.MustListLinks(),
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"maxReadPosts":    models.MustListMaxReadPost(),
		"maxCommentPosts": models.MustListMaxCommentPost(),
	})
}
