package models

type Comment struct {
	Model

	Content string `form:"content"`
	UserID uint64
	PostID uint64
	User      User   `binding:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
}