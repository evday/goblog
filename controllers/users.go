package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
	"myblog/models"
)

func UserIndex(c *gin.Context) {
	db := models.GetDB()
	var users []models.User
	db.Find(&users)
	h := DefaultH(c)
	h["Title"] = "用户列表"
	h["Users"] = users
	c.HTML(http.StatusOK,"users/index",h)
}

func UserDelete(c *gin.Context) {
	db := models.GetDB()
	user := models.User{}
	db.First(&user,c.Param("id"))
	if user.ID == 0 {
		c.HTML(http.StatusNotFound,"errors/404",nil)
		return
	}
	if err := db.Delete(&user).Error;err != nil {
		c.HTML(http.StatusInternalServerError,"errors/500",nil)
		logrus.Error(err)
		return	
	}
	c.Redirect(http.StatusFound,"/admin/users")
}