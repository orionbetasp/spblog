package controllers

import (
	"math"
	"net/http"
	"spblog/conf"
	"spblog/models"
	"spblog/util"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func TagCreate(c *gin.Context) {
	name := c.PostForm("name")
	tag := &models.Tag{Name: name}
	util.Logger.Info("create", zap.Any("tag", tag))
	err := tag.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "tag/new.html", gin.H{
			"message": err.Error(),
			"tag":     tag,
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/tag")
}

func TagGet(c *gin.Context) {
	var (
		tagName   string
		page      string
		pageIndex int
		pageSize  = conf.Con.GetInt("page_size")
		total     int
		err       error
		policy    *bluemonday.Policy
		posts     []*models.Post
	)
	tagName = c.Param("tag")
	page = c.Query("tag")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	posts, err = models.ListPublishedPost(tagName, pageIndex, pageSize)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByTag(tagName)
	if err != nil {
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

func TagNew(c *gin.Context) {
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "tag/new.html", gin.H{
		"user": user,
	})
}

func TagEdit(c *gin.Context) {
	id := c.Param("id")
	tag, err := models.GetTagById(id)
	if err != nil {
		Handle404(c)
	}
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "tag/modify.html", gin.H{
		"tag":  tag,
		"user": user,
	})
}

func TagUpdate(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	tag := &models.Tag{Name: name}
	tag.ID = uint(pid)
	util.Logger.Info("update", zap.Any("tag", tag))
	err = tag.Update()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/tag")
}

func TagDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	tag := &models.Tag{}
	tag.ID = uint(pid)
	util.Logger.Info("delete", zap.Any("tag", tag))
	err = tag.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func TagIndex(c *gin.Context) {
	tags, _ := models.ListAllTag()
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "admin/tag.html", gin.H{
		"tags":     tags,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}
