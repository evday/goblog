package controllers

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)
//返回404页面
func NotFound(c *gin.Context){
	ShowErrorPage(c,http.StatusNotFound,nil)
}

//方法不被允许
func MethodNotAllowed(c *gin.Context){
	ShowErrorPage(c,http.StatusMethodNotAllowed,nil)
}

//返回错误页面
func ShowErrorPage(c *gin.Context,code int,err error){
	H := DefaultH(c)
	H["Error"] = err
	c.HTML(code,fmt.Sprintf("errors/%d",code),H)
}