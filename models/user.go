package models

import (
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Register struct {
	Name string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
	Avatar string `form:"avatar"`
}

type User struct {
	Model

	Email    string `form:"email" binding:"required"`
	Name     string `form:"name"`
	Image string `form:"image"`
	Password string `form:"password" binding:"required"`
}

//保存之前给密码加密
func (u *User) BeforeSave() (err error) {
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hash)
	return
}