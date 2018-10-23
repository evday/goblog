package controllers

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/Sirupsen/logrus"
	"myblog/models"
	"strings"


)
func PostGet(c *gin.Context){
	db := models.GetDB()
	session := sessions.Default(c)
	post := models.Post{}
	
	db.Preload("Tags").Preload("User").Preload("Comments",func(db *gorm.DB)*gorm.DB{
		return db.Order("comments.created_at DESC")
	}).First(&post,c.Param("id"))

	if post.ID ==0 || !post.Published {
		c.HTML(http.StatusNotFound,"errors/404",nil)
		return
	}
	post.View++
	post.UpdateView()

	h := DefaultH(c)
	h["Title"] = post.Title
	h["Post"] = post
	h["Flash"] = session.Flashes()
	prepost := models.Post{}
	nextpost := models.Post{}
	db.Select("id,title").Order("id desc").Where("id > ?",c.Param("id")).Find(&prepost).Limit(1)
	db.Select("id,title").Order("id desc").Where("id < ?",c.Param("id")).Find(&nextpost).Limit(1)
	h["PrePost"] = prepost
	h["NextPost"] = nextpost

	tag := c.Param("tag")
	if tag != "" {
		h["Tag"] = tag
	}else{
		h["Tag"] = post.Tags[0].Title
	}
	
	session.Save()
	c.HTML(http.StatusOK,"posts/show",h)
}

func PostIndex(c *gin.Context) {
	db := models.GetDB()
	var posts []models.Post
	db.Preload("Tags").Find(&posts)
	h := DefaultH(c)
	h["Title"] = "文章列表"
	h["Posts"] = posts
	c.HTML(http.StatusOK, "admin/index", h)
}

func PostNew(c *gin.Context) {
	var tags []models.Tag
	db := models.GetDB()
	db.Order("title asc").Find(&tags)
	h := DefaultH(c)
	h["Title"] = "添加文章"
	h["Tags"] = tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "admin/form", h)
}

func PostCreate(c *gin.Context) {
	session := sessions.Default(c)
	post := models.Post{}
	db := models.GetDB()

	if err := c.ShouldBind(&post);err != nil {
		
		session.AddFlash(err.Error())
		session.Save()
		logrus.Error(err)
		c.Redirect(http.StatusSeeOther,"/admin/new_post")
		return
	}
	uri,err := UploadPost(c)
	if err != nil {
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther,"/admin/new_post")
		return
	}
	tags := make([]models.Tag,0,len(post.FormTags))
	for i := range post.FormTags{
		tags = append(tags,models.Tag{Title:post.FormTags[i]})
	}
	post.Image = uri
	post.Tags = tags
	post.CreatedAt = post.CreatedAt
	if user,exists := c.Get("User");exists {
		post.UserID = user.(*models.User).ID
	}
	if err := db.Create(&post).Error;err != nil {
		c.HTML(http.StatusInternalServerError,"errors/500",nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound,"/admin")

}

func PostEdit(c *gin.Context){
	db := models.GetDB()
	post := models.Post{}
	db.Preload("Tags").First(&post,c.Param("id"))
	if post.ID == 0 {
		c.HTML(http.StatusNotFound,"errors/404",nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "编辑文章"
	h["Post"] = post
	h["Tags"] = post.Tags
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK,"admin/form",h)
}

func PostUpdate(c *gin.Context){
	_, header , _ := c.Request.FormFile("image")
	filename := header.Filename

	
	
	
	session := sessions.Default(c)
	db := models.GetDB()
	post := models.Post{}
	if err := c.ShouldBind(&post);err != nil {
		
		session.AddFlash(err.Error())
		session.Save()
		logrus.Error(err)
		c.Redirect(http.StatusSeeOther,fmt.Sprintf("/admin/posts/%s/edit",c.Param("id")))
		return
	}
	content := c.PostForm("content")

	db.First(&post,c.Param("id"))
	if post.Image == "" || !strings.Contains(post.Image,filename) {
		uri,err := UploadPost(c)
		if err != nil {
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(http.StatusSeeOther,"/admin/new_post")
			return
		}
		post.Image = uri
	}

	tags := make([]models.Tag,0,len(post.FormTags))
	for i := range post.FormTags{
		tags = append(tags,models.Tag{Title:post.FormTags[i]})
	}
	post.Content = content
	post.Tags = tags

	if err := db.Save(&post).Error;err != nil {
		c.HTML(http.StatusInternalServerError,"errors/500",nil)
		logrus.Error(err)
		return
	}
	if err := db.Exec("DELETE FROM posts_tags WHERE post_id = ? AND tag_title NOT IN(?)",post.ID,post.FormTags).Error;err != nil {
		c.HTML(http.StatusInternalServerError,"errors/500",nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound,"/admin")
}

func PostDelete(c *gin.Context){
	db := models.GetDB()
	post := models.Post{}
	db.First(&post,c.Param("id"))
	if post.ID == 0 {
		c.HTML(http.StatusNotFound,"errors/404",nil)
		return
	}
	if err := db.Delete(&post).Error;err != nil {
		c.HTML(http.StatusInternalServerError,"errors/500",nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound,"/admin")
}


func SearchPost(c *gin.Context){
	keyword := c.PostForm("keyboard")
	db := models.GetDB()
	var list []models.Post
	db.Preload("Tags").Preload("User").Where("Title LIKE ?",fmt.Sprintf("%%%s%%",keyword)).Where("published = true").Order("id desc").Find(&list)
	if len(list) == 0  {
		c.HTML(http.StatusNotFound,"posts/404",nil)
		return
	}
	h := DefaultH(c)
	h["Posts"] = list
	h["KeyWord"] = keyword
	c.HTML(http.StatusOK,"posts/search",h)

}