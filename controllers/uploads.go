package controllers

import (
	"encoding/hex"
	"crypto/md5"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	// "fmt"
	// "time"
	"myblog/system"
	"github.com/Sirupsen/logrus"

)

func UploadPost(c *gin.Context) (string, error){
	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	mpartFile, mpartHeader, err := c.Request.FormFile("image")
	if mpartHeader.Filename == ""{
		return "/public/images/zd02.jpg",nil
	}
	if err != nil {
		logrus.Error(err)
		c.String(400, err.Error())
		return "", err
	}
	defer mpartFile.Close()
	uri, err := saveFile(mpartHeader, mpartFile)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return uri, nil
}

func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := filepath.Ext(fh.Filename)
	// newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name ;D
	newName := Md5(fh.Filename) + fileExt //unique file name ;D
	uri := "/public/uploads/" + newName
	fullName := filepath.Join(system.UploadsPath(), newName)

	file, err := os.OpenFile(fullName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}
	return uri, nil
}

func Md5(filename string)(md5Str string){
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(filename))
	cipherStr := md5Ctx.Sum(nil)
	md5Str = hex.EncodeToString(cipherStr)
	return

}