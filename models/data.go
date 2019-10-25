package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// I don't need soft delete,so I use customized BaseModel instead gorm.Model
type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// table posts
type Post struct {
	BaseModel
	Title        string     // title
	Body         string     `gorm:"size:9999"` // body
	View         int        // view count
	IsPublished  bool       // published or not
	Tags         []*Tag     `gorm:"-"` // tags of post
	Comments     []*Comment `gorm:"-"` // comments of post
	CommentTotal int        `gorm:"-"` // count of comment
}

// WxAppData
type WxAppData struct {
	BaseModel
	Name string `json:"name"`                  // user name
	Data string `gorm:"size:9999" json:"list"` // user data
}

// table tags
type Tag struct {
	BaseModel
	Name  string // tag name
	Total int    `gorm:"-"` // count of post
}

// table post_tags
type PostTag struct {
	BaseModel
	PostId uint // post id
	TagId  uint // tag id
}

// table users
type User struct {
	gorm.Model
	Email       string `gorm:"unique_index;default:null"` //邮箱
	Telephone   string `gorm:"unique_index;default:null"` //手机号码
	Password    string `gorm:"default:null"`              //密码
	VerifyState string `gorm:"default:'0'"`               //邮箱验证状态
	SecretKey   string `gorm:"default:null"`              //密钥
	//OutTime       time.Time //过期时间
	GithubLoginId string `gorm:"unique_index;default:null"` // github唯一标识
	GithubUrl     string //github地址
	IsAdmin       bool   //是否是管理员
	AvatarUrl     string // 头像链接
	NickName      string // 昵称
	LockState     bool   `gorm:"default:'0'"` //锁定状态
}

// table comments
type Comment struct {
	BaseModel
	UserID    uint   // 用户id
	Content   string // 内容
	PostID    uint   // 文章id
	ReadState bool   `gorm:"default:'0'"` // 阅读状态
	//Replies []*Comment // 评论
	NickName  string `gorm:"-"`
	AvatarUrl string `gorm:"-"`
	GithubUrl string `gorm:"-"`
}

// table subscribe
type Subscriber struct {
	gorm.Model
	Email          string    `gorm:"unique_index"` //邮箱
	VerifyState    bool      `gorm:"default:'0'"`  //验证状态
	SubscribeState bool      `gorm:"default:'1'"`  //订阅状态
	OutTime        time.Time //过期时间
	SecretKey      string    // 秘钥
	Signature      string    //签名
}

// table link
type Link struct {
	gorm.Model
	Name string //名称
	Url  string //地址
	Sort int    `gorm:"default:'0'"` //排序
	View int    //访问次数
}

// query result
type QrArchive struct {
	ArchiveDate time.Time //month
	Total       int       //total
	Year        int       // year
	Month       int       // month
}
