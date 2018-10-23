package models

type Post struct {
	Model

	Title string `form:"title" binding:"required"`
	Content string `form:"content" gorm:"type:text"`
	Image string `form:"image"`
	Published bool `form:"published"`
	UserID uint64
	View   int
	User User `binding:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
	FormTags []string `form:"tags" gorm:"-"`
	Tags []Tag `binding:"-" form:"-" json:"tags" gorm:"many2many:posts_tags;"`
	Comments  []Comment `binding:"-"`
}


func (post *Post) UpdateView() error {
	db := GetDB()
	return db.Model(post).Updates(map[string]interface{}{
		"view":post.View,
	}).Error
}