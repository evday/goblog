package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
	"myblog/system"
	"myblog/models"
	"myblog/router"
)

func main() {
	//初始化日志函数
	initLogger()
	//处理配置文件
	system.LoadConfig()
	//函数GetConnectionString获得数据库信息，SetDB设置数据库连接
	models.SetDB(system.GetConnectionString())
	//自动初始化数据库，使用的是gorm这个orm
	models.AutoMigrate()
	//静态文件
	system.LoadTemplates()

	router := router.InitRouter()
	router.Run(":8080")	
}

//打印日志
func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}