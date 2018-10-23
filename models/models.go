package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


type Model struct {
	ID        uint64     `form:"id" gorm:"primary_key"`
	CreatedAt time.Time  `binding:"-" form:"-"`
	UpdatedAt time.Time  `binding:"-" form:"-"`
	DeletedAt *time.Time `binding:"-" form:"-"`
}

var db *gorm.DB

//设置数据库连接
func SetDB(connection string) {
	var err error
	db, err = gorm.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
}


func GetDB() *gorm.DB {
	return db
}


//自动提交
func AutoMigrate() {
	db.AutoMigrate(&User{}, &Tag{},&Post{}, &Comment{})
}


