package controllers

import (
	"encoding/json"
	"spblog/models"
	"spblog/util"

	"github.com/gin-gonic/gin"
	//"github.com/tidwall/gjson"
)

func WxAppUpData(c *gin.Context) {
	var wa models.WxAppData
	wa.Name = c.Query("name")
	wa.Data = c.Query("list")
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
	var wa models.WxAppData
	wa.Name = c.Query("name")
	if wa.Name == "" {
		c.String(200, "不科学，没名字！")
		return
	}
	waRes, err := models.GetWxAppDataByName(wa.Name)
	if err != nil {
		c.String(200, "就出问题了呗！")
		return
	}
	if waRes.Name == "" {
		c.String(200, "你之前有没有上传数据，心里没点B数嘛！")
		return
	}
	j, err := json.Marshal(waRes)
	if err != nil || wa.Name == "" {
		c.String(200, "就出问题了呗！")
		return
	}
	util.Logger.Info("下载数据：" + string(j))
	c.String(200, string(j))
}
