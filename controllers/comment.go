package controllers

import (
	"fmt"
	"spblog/conf"
	"spblog/models"
	"spblog/util"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CommentPost(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
	)
	defer writeJSON(c, res)
	s := sessions.Default(c)
	sessionUserID := s.Get(SessionKey)
	userId, _ := sessionUserID.(uint)

	verifyCode := c.PostForm("verifyCode")
	captchaId := s.Get(SessionCaptcha)
	s.Delete(SessionCaptcha)
	_captchaId, _ := captchaId.(string)
	if !captcha.VerifyString(_captchaId, verifyCode) {
		res["message"] = "error verifycode,验证码错误呢！"
		return
	}

	postId := c.PostForm("postId")
	content := c.PostForm("content")
	if len(content) == 0 {
		res["message"] = "content cannot be empty，空的评论啥呢！"
		return
	}

	post, err = models.GetPostById(postId)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	pid, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := &models.Comment{
		PostID:  uint(pid),
		Content: content,
		UserID:  userId,
	}
	err = comment.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = NotifyEmail("[spblog]您有一条新评论", fmt.Sprintf("<a href=\"%s/post/%d\" target=\"_blank\">%s</a>:%s",
		conf.Con.GetString("domain"), post.ID, post.Title, content))
	if err != nil {
		util.Logger.Error(err.Error())
	}
	res["succeed"] = true
}

func CommentDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
		cid uint64
	)
	defer writeJSON(c, res)

	s := sessions.Default(c)
	sessionUserID := s.Get(SessionKey)
	userId, _ := sessionUserID.(uint)

	commentId := c.Param("id")
	cid, err = strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := &models.Comment{
		UserID: userId,
	}
	comment.ID = uint(cid)
	err = comment.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func CommentRead(c *gin.Context) {
	var (
		id  string
		_id uint64
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id = c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := new(models.Comment)
	comment.ID = uint(_id)
	err = comment.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func CommentReadAll(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	err = models.SetAllCommentRead()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
