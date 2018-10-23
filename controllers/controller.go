package controllers

import (
	"net/http"
	csrf "github.com/utrack/gin-csrf"
	"github.com/gin-gonic/gin"
)

const userIDkey = "UserID"

func DefaultH(c *gin.Context)gin.H{
	return gin.H{
		"Title":"",
		"Context":c,
		"Csrf":csrf.GetToken(c),
	}
}

func MarkDownTest(c *gin.Context){
	c.HTML(http.StatusOK,"home/test",nil)
}