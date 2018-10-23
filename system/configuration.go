package system

import (
	"os"
	"path"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
)

type Configs struct {
	Debug Config
	Release Config
}

type Config struct {
	Public string 	`json:"public"`
	Domain string `json:"domain"`
	SessionSecret string `json:"session_secret"`
	SignupEnabled bool `json:"signup_enabled"`
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Name string
	User string
	Password string
	IP string
	Port string
}

var config *Config
//配置文件
func LoadConfig() {
	data,err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	configs := &Configs{}
	err = json.Unmarshal(data,configs)
	if err != nil {
		panic(err)
	}
	switch gin.Mode() {
	case gin.DebugMode:
		config = &configs.Debug
	case gin.ReleaseMode:
		config = &configs.Release
	default:
		panic(fmt.Sprintf("Unknown gin mode %s",gin.Mode()))
	}

	if !path.IsAbs(config.Public){
		workingDir,err := os.Getwd()
		if err != nil {
			panic(err)
		}
		config.Public = path.Join(workingDir,config.Public)
	}
}

//返回配置文件对象
func GetConfig() *Config{
	return config
}

//静态文件地址
func PublicPath() string {
	return config.Public
}
//返回文件上传地址
func UploadsPath() string {
	return path.Join(config.Public,"uploads")
}

//返回数据库配置信息
func GetConnectionString()string{
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",config.Database.User,config.Database.Password,config.Database.IP,config.Database.Port,config.Database.Name)
}