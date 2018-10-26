package controllers

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"myblog/models"
)

func HomeGet(c *gin.Context){
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	db := models.GetDB()
	
	var list []models.Post
	

	paginator := Pagging(&Param{
		DB:db,
		Page:page,
		Limit:10,
		OrderBy:[]string{"id desc"},
		ShowSQL:false,
	},&list)

	// c.JSON(200,paginator)

	
	h := DefaultH(c)
	h["Title"] = "Welcome to basic GIN blog"
	h["pagination"] = paginator
	c.HTML(http.StatusOK,"home/index",h)
}


func AboutGet(c *gin.Context){
	c.HTML(http.StatusOK,"home/about",nil)
}