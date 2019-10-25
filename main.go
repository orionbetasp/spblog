package main

import (
	"spblog/conf"
	"spblog/models"
	"spblog/router"
	"spblog/util"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	defer func() { _ = db.Close() }()
	util.Logger.Info("InitDB success")
	r := router.SetupRouter()
	util.Logger.Info("blog server start")
	err = r.Run(conf.Con.GetString("addr"))
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
}
