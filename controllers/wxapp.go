package controllers

import (
	"encoding/base64"
	"encoding/json"
	"spblog/models"
	"spblog/util"

	"github.com/gin-gonic/gin"
	//"github.com/tidwall/gjson"
)

func WxAppUpData(c *gin.Context) {
	var wa models.WxAppData
	n, _ := base64.URLEncoding.DecodeString(c.Query("name"))
	d, _ := base64.URLEncoding.DecodeString(c.Query("list"))
	wa.Name = string(n)
	wa.Data = string(d)
	util.Logger.Info(string(n) + " 上传数据：" + string(d))
	if wa.Name == "" {
		c.String(200, "不科学，没名字！")
		return
	}
	err := wa.Insert()
	if err != nil {
		c.String(200, "就出问题了呗！"+err.Error())
		return
	}
	c.String(200, wa.Name+"の数据已同步!")
}

func WxAppDownData(c *gin.Context) {
	n, _ := base64.URLEncoding.DecodeString(c.Query("name"))
	util.Logger.Info(string(n) + " 下载数据：")
	if string(n) == "" {
		c.String(200, "不科学，没名字！")
		return
	}
	waRes, err := models.GetWxAppDataByName(string(n))
	if err != nil {
		c.String(200, "就出问题了呗！")
		return
	}
	if waRes.Name == "" {
		c.String(200, "你之前有没有上传数据，心里没点B数嘛！")
		return
	}
	j, err := json.Marshal(waRes)
	if err != nil {
		c.String(200, "就出问题了呗！"+err.Error())
		return
	}
	util.Logger.Info("下载数据：" + string(j))
	c.String(200, string(j))
}
