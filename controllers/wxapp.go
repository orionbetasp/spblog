package controllers

import (
	"encoding/json"
	"io/ioutil"
	"spblog/models"

	"github.com/gin-gonic/gin"
	//"github.com/tidwall/gjson"
)

func WxAppUpData(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	//name := gjson.Get(string(body), "name").String()
	//data := gjson.Get(string(body), "list").String()
	var wa models.WxAppData
	err := json.Unmarshal(body, &wa)
	if err != nil || wa.Name == "" {
		c.String(200, "就出问题了呗！")
		return
	}
	err = wa.Insert()
	if err != nil {
		c.String(200, "就出问题了呗！")
		return
	}
	c.String(200, wa.Name+"の数据已同步!")
}

func WxAppDownData(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	//name := gjson.Get(string(body), "name").String()
	//data := gjson.Get(string(body), "list").String()
	var wa models.WxAppData
	err := json.Unmarshal(body, &wa)
	if err != nil || wa.Name == "" {
		c.String(200, "就出问题了呗！")
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
	c.String(200, string(j))
}
