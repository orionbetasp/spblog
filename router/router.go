package router

import (
	"html/template"
	"os"
	"path/filepath"
	"spblog/conf"
	"spblog/controllers"
	"spblog/models"
	"spblog/util"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/claudiu/gocron"
	"github.com/gin-gonic/gin"
)

//SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	setTemplate(r)
	setSessions(r)

	r.Use(gin.Recovery())
	r.Use(SharedData())
	r.Use(signAuthMiddleware())

	//Periodic tasks
	gocron.Every(1).Day().Do(controllers.CreateXMLSitemap)
	gocron.Start()

	r.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	r.StaticFile("/favicon.ico", filepath.Join(getCurrentDirectory(), "./static/favicon.ico"))

	r.NoRoute(controllers.Handle404)

	//小程序
	r.GET("/UpData", controllers.WxAppUpData)
	r.GET("/DownData", controllers.WxAppDownData)

	r.GET("/", controllers.IndexGet)
	r.GET("/index", controllers.IndexGet)
	r.GET("/rss", controllers.RssGet)

	if conf.Con.GetBool("signup_enabled") {
		r.GET("/signup", controllers.SignupGet)
		r.POST("/signup", controllers.SignupPost)
	}
	// user signin and logout
	r.GET("/signin", controllers.SigninGet)
	r.POST("/signin", controllers.SigninPost)
	r.GET("/logout", controllers.LogoutGet)
	r.GET("/oauth2callback", controllers.Oauth2Callback)
	r.GET("/auth/:authType", controllers.AuthGet)

	// captcha
	r.GET("/captcha", controllers.CaptchaGet)

	visitor := r.Group("/visitor")
	visitor.Use(AuthRequired())
	{
		visitor.POST("/new_comment", controllers.CommentPost)
		visitor.POST("/comment/:id/delete", controllers.CommentDelete)
	}

	// subscriber
	r.GET("/subscribe", controllers.SubscribeGet)
	r.POST("/subscribe", controllers.Subscribe)
	r.GET("/active", controllers.ActiveSubscriber)
	r.GET("/unsubscribe", controllers.UnSubscribe)

	//r.GET("/tag/:id", controllers.PageGet)
	r.GET("/post/:id", controllers.PostGet)
	r.GET("/tag/:tag", controllers.TagGet)
	r.GET("/archives/:year/:month", controllers.ArchiveGet)

	r.GET("/link/:id", controllers.LinkGet)

	authorized := r.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		// index
		authorized.GET("/index", controllers.AdminIndex)

		// tag
		authorized.GET("/tag", controllers.TagIndex)
		authorized.GET("/new_tag", controllers.TagNew)
		authorized.POST("/new_tag", controllers.TagCreate)
		authorized.GET("/tag/:id/edit", controllers.TagEdit)
		authorized.POST("/tag/:id/edit", controllers.TagUpdate)
		//authorized.POST("/tag/:id/publish", controllers.TagPublish)
		authorized.POST("/tag/:id/delete", controllers.TagDelete)

		// post
		authorized.GET("/post", controllers.PostIndex)
		authorized.GET("/new_post", controllers.PostNew)
		authorized.POST("/new_post", controllers.PostCreate)
		authorized.GET("/post/:id/edit", controllers.PostEdit)
		authorized.POST("/post/:id/edit", controllers.PostUpdate)
		authorized.POST("/post/:id/publish", controllers.PostPublish)
		authorized.POST("/post/:id/delete", controllers.PostDelete)

		// user
		authorized.GET("/user", controllers.UserIndex)
		authorized.POST("/user/:id/lock", controllers.UserLock)

		// profile
		authorized.GET("/profile", controllers.ProfileGet)
		authorized.POST("/profile", controllers.ProfileUpdate)
		authorized.POST("/profile/email/bind", controllers.BindEmail)
		authorized.POST("/profile/email/unbind", controllers.UnbindEmail)
		authorized.POST("/profile/github/unbind", controllers.UnbindGithub)

		// subscriber
		authorized.GET("/subscriber", controllers.SubscriberIndex)
		authorized.POST("/subscriber", controllers.SubscriberPost)

		// link
		authorized.GET("/link", controllers.LinkIndex)
		authorized.POST("/new_link", controllers.LinkCreate)
		authorized.POST("/link/:id/edit", controllers.LinkUpdate)
		authorized.POST("/link/:id/delete", controllers.LinkDelete)

		// comment
		authorized.POST("/comment/:id", controllers.CommentRead)
		authorized.POST("/read_all", controllers.CommentReadAll)

		// mail
		authorized.POST("/new_mail", controllers.SendMail)
		authorized.POST("/new_batchmail", controllers.SendBatchMail)
	}
	return r
}

func setTemplate(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": util.DateFormat,
		"substring":  util.Substring,
		"isOdd":      util.IsOdd,
		"isEven":     util.IsEven,
		"truncate":   util.Truncate,
		"add":        util.Add,
		"minus":      util.Minus,
		"listtag":    ListTag,
	}

	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
}

func ListTag() (tagstr string) {
	tags, err := models.ListTag()
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	tagNames := make([]string, 0)
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	tagstr = strings.Join(tagNames, ",")
	return
}

//setSessions initializes sessions & csrf middlewares
func setSessions(router *gin.Engine) {
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := cookie.NewStore([]byte(conf.Con.GetString("session_secret")))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
	//https://github.com/utrack/gin-csrf
	/*router.Use(csrf.Middleware(csrf.Options{
		Secret: config.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))*/
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		util.Logger.Error(err.Error())
	}
	util.Logger.Info(dir)
	//return strings.Replace(dir, "\\", "/", -1)
	return ""
}
