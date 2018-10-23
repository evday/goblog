package controllers

import (
	"net/url"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"myblog/models"
	"myblog/system"
)

func ContextData() gin.HandlerFunc {
	return func(c *gin.Context){
		session := sessions.Default(c)
		if uID := session.Get(userIDkey);uID != nil {
			user := models.User{}
			models.GetDB().First(&user,uID)
			if user.ID !=0 {
				c.Set("User",&user)
			}
		}

		if system.GetConfig().SignupEnabled {
			c.Set("SignupEnabled",true)
		}

		c.Next()
	}
}


func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context){
		if user,_ := c.Get("User");user != nil {
			c.Next()
		}else{
			c.Redirect(http.StatusFound,fmt.Sprintf("/login?return=%s",url.QueryEscape(c.Request.RequestURI)))
			c.Abort()
		}
	}
}