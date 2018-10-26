package router

import (
	"github.com/gin-gonic/gin"
	"myblog/system"
	"myblog/controllers"
	"net/http"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions"
	"github.com/utrack/gin-csrf"
	"github.com/Sirupsen/logrus"
	"fmt"
)

func InitRouter() *gin.Engine{
	router := gin.Default()
	router.SetHTMLTemplate(system.GetTemplates())
	router.StaticFS("/public", http.Dir(system.PublicPath()))

	//设置session，从配置文件中取出SessionSecret
	config := system.GetConfig()
	store := memstore.NewStore([]byte(config.SessionSecret))
	//设置session过期时间
	store.Options(sessions.Options{HttpOnly:true,MaxAge:7*86400})
	router.Use(sessions.Sessions("gin-session",store))
	//设置csrftoken
	router.Use(csrf.Middleware(csrf.Options{
		Secret:config.SessionSecret,
		ErrorFunc:func(c *gin.Context){
			logrus.Error("CSRF token mismatch")
			controllers.ShowErrorPage(c,400,fmt.Errorf("CSRF token mismatch"))
			c.Abort()
		},
	}))

	router.Use(controllers.ContextData())
	
	//主页
	router.GET("/", controllers.HomeGet)

	//用户登录
	router.GET("/login",controllers.SignInGet)
	router.GET("/about",controllers.AboutGet)
	router.POST("/login",controllers.SignInPost)
	router.GET("/register",controllers.RegisterGet)
	router.POST("/register",controllers.RegisterPost)
	router.GET("/logout", controllers.LogoutGet)

	router.GET("/posts/:id/:tag", controllers.PostGet)
	router.GET("/posts/:id", controllers.PostGet)
	router.GET("/tags/:title",controllers.TagGet)

	router.POST("/search",controllers.SearchPost)

	admin := router.Group("/admin")
	admin.Use(controllers.AuthRequired())
	{
		admin.GET("/",controllers.PostIndex)
		admin.GET("/new_post", controllers.PostNew)
		admin.POST("/new_post", controllers.PostCreate)
		admin.GET("/posts/:id/edit",controllers.PostEdit)
		admin.POST("/posts/:id/edit",controllers.PostUpdate)
		admin.POST("/posts/:id/delete",controllers.PostDelete)

		admin.GET("/tags", controllers.TagIndex)
		admin.GET("/new_tag", controllers.TagNew)
		admin.POST("/new_tag", controllers.TagCreate)
		admin.POST("/tags/:title/delete", controllers.TagDelete)

		admin.GET("/users",controllers.UserIndex)
		admin.POST("/users/:id/delete",controllers.UserDelete)
		
	}
	

	return router


}