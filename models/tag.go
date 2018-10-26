package models



//标签
type Tag struct {


	Title string `binding:"required" form:"title" gorm:"primary_key"`
	Posts []Post `gorm:"many2many:posts_tags;foreignkey:title"`
}