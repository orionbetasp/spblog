package controllers

import (
	"net/http"
	"spblog/util"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"spblog/models"

	"github.com/gin-gonic/gin"
)

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil || !post.IsPublished {
		Handle404(c)
		return
	}
	post.View++
	_ = post.UpdateView()
	post.Tags, _ = models.ListTagByPostId(id)
	post.Comments, _ = models.ListCommentByPostID(id)
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "post/display.html", gin.H{
		"post": post,
		"user": user,
	})
}

func PostNew(c *gin.Context) {
	tags, _ := models.ListAllTag()
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "post/new.html", gin.H{
		"inTags": tags,
		"user":   user,
	})
}

func PostCreate(c *gin.Context) {
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished

	post := &models.Post{
		Title:       title,
		Body:        body,
		IsPublished: published,
	}
	util.Logger.Info("create", zap.Any("post", post))
	err := post.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
		return
	}

	// add tag for post
	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tag := range tagArr {
			tagId, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				continue
			}
			pt := &models.PostTag{
				PostId: post.ID,
				TagId:  uint(tagId),
			}
			_ = pt.Insert()
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/post")
}

func PostEdit(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		Handle404(c)
		return
	}
	post.Tags, _ = models.ListTagByPostId(id)
	tags, _ := models.ListAllTag()
	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "post/modify.html", gin.H{
		"post":   post,
		"inTags": tags,
		"user":   user,
	})
}

func PostUpdate(c *gin.Context) {
	id := c.Param("id")
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished

	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		Handle404(c)
		return
	}

	post := &models.Post{
		Title:       title,
		Body:        body,
		IsPublished: published,
	}
	util.Logger.Info("update", zap.Any("post", post))
	post.ID = uint(pid)
	err = post.Update()
	if err != nil {
		c.HTML(http.StatusOK, "post/modify.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
		return
	}
	// 删除tag
	_ = models.DeletePostTagByPostId(post.ID)
	// 添加tag
	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tag := range tagArr {
			tagId, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				continue
			}
			pt := &models.PostTag{
				PostId: post.ID,
				TagId:  uint(tagId),
			}
			_ = pt.Insert()
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/post")
}

func PostPublish(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	post, err = models.GetPostById(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	post.IsPublished = !post.IsPublished
	err = post.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func PostDelete(c *gin.Context) {
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
	post := &models.Post{}
	post.ID = uint(pid)
	util.Logger.Info("delete", zap.Any("post", post))
	err = post.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	_ = models.DeletePostTagByPostId(uint(pid))
	res["succeed"] = true
}

func PostIndex(c *gin.Context) {
	posts, _ := models.ListAllPost("")
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.Itoa(int(post.ID)))
	}

	user, _ := c.Get(ContextUserKey)
	c.HTML(http.StatusOK, "admin/post.html", gin.H{
		"posts":    posts,
		"Active":   "posts",
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}
