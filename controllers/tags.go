package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/Sirupsen/logrus"
	"myblog/models"
)

func TagGet(c *gin.Context) {
	db := models.GetDB()
	tag := models.Tag{}
	db.Preload("Posts","published = true").Preload("Posts.Comments").Preload("Posts.Tags").Preload("Posts.User").Find(&tag,"title = ? ",c.Param("title"))
	if len(tag.Title) == 0 {
		c.HTML(http.StatusNotFound,"errors/404",nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = tag.Title
	h["Tag"] = tag
	c.HTML(http.StatusOK,"tags/show",h)
}

func TagIndex(c *gin.Context){
	db := models.GetDB()
	var tags []models.Tag
	db.Preload("Posts").Order("title asc").Find(&tags)
	h := DefaultH(c)
	h["Title"] = "标签列表"
	h["Tags"] = tags
	c.HTML(http.StatusOK,"tags/index",h)

}


func TagNew(c *gin.Context){
	h := DefaultH(c)
	h["Title"] = "添加标签"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK,"tags/form",h)
}

func TagCreate(c *gin.Context){
	tag := models.Tag{}
	db := models.GetDB()
	if err := c.ShouldBind(&tag);err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther,"/admin/new_tag")
		return
	}

	if err := db.Create(&tag).Error;err != nil {
		c.HTML(http.StatusInternalServerError,"errors/500",nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound,"/admin/tags")
}

func TagDelete(c *gin.Context){
	db := models.GetDB()
	tag := models.Tag{}
	db.Where("title = ?", c.Param("title")).First(&tag)
	
	if len(tag.Title) == 0{
		c.HTML(http.StatusNotFound,"errors/404",nil)
		return 
	}

	if err := db.Delete(&tag).Error;err != nil {
		logrus.Error(err)
		c.HTML(http.StatusInternalServerError,"errors/500",gin.H{
			"Error":err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound,"/admin/tags")
}