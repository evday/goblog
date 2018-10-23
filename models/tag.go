package models

import "time"

//标签
type Tag struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Title string `binding:"required" form:"title" gorm:"primary_key"`
	Posts []Post `gorm:"many2many:posts_tags;foreignkey:title"`
}