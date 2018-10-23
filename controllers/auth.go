package controllers

import (
	
	"fmt"
	"net/http"
	"net/url"

	"strings"

	"github.com/Sirupsen/logrus"
	"myblog/models"
	
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//获取登录页面
func SignInGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Basic GIN web-site signin form"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/login", h)
}

//登录
func SignInPost(c *gin.Context) {
	session := sessions.Default(c)
	login := models.Login{}
	db := models.GetDB()
	returnURL := c.DefaultQuery("return", "/admin")
	if err := c.ShouldBind(&login); err != nil {
		session.AddFlash("Please, fill out form correctly.")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User

	db.Where("email = lower(?)", login.Email).First(&user)
	
	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), login.Email)
		session.AddFlash("Email or password incorrect")
		session.Save()
		c.Redirect(http.StatusFound, fmt.Sprintf("/login?return=%s", url.QueryEscape(returnURL)))
		return
	}

	session.Set(userIDkey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, returnURL)
}

//注册页面
func RegisterGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Basic GIN web-site signup form"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/register", h)
}

//用户注册
func RegisterPost(c *gin.Context) {
	session := sessions.Default(c)
	register := models.Register{}
	db := models.GetDB()
	if err := c.ShouldBind(&register); err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/register")
		return
	}
	register.Email = strings.ToLower(register.Email)
	user := models.User{}
	db.Where("email = ?", register.Email).First(&user)
	if user.ID != 0 {
		session.AddFlash("User exists")
		session.Save()
		c.Redirect(http.StatusFound, "/register")
		return
	}

	uri,err := UploadPost(c)
	if err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/register")
		return
	}
	
	//创建用户
	user.Email = register.Email
	user.Password = register.Password
	user.Name = register.Name
	user.Image = uri

	if err := db.Create(&user).Error; err != nil {
		session.AddFlash("Error whilst registering user.")
		session.Save()
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/register")
		return
	}

	session.Set(userIDkey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, "/")
	return
}

//用户登出
func LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userIDkey)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}
