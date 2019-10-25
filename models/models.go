package models

import (
	"spblog/util"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:root@/spblog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		util.Logger.Error(err.Error())
		return db, err
	}
	//db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Minute)

	db.AutoMigrate(&WxAppData{}, &Post{}, &Tag{}, &PostTag{}, &User{}, &Comment{}, &Subscriber{}, &Link{})
	db.Model(&PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")

	//db, err := gorm.Open("sqlite3", conf.Con.GetString("dns"))
	////db, err := gorm.Open("mysql", "root:mysql@/wblog?charset=utf8&parseTime=True&loc=Asia/Shanghai")
	//if err != nil {
	//	return db, err
	//}
	//
	//DB = db
	////db.LogMode(true)
	//db.AutoMigrate(&Page{}, &Post{}, &Tag{}, &PostTag{}, &User{}, &Comment{}, &Subscriber{}, &Link{}, &SmmsFile{})
	//db.Model(&PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")
	DB = db
	return db, nil
}
