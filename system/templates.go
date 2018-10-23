package system

import (
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
	"fmt"
	"time"
	"os"
	"strings"
	"html/template"
	"path/filepath"
	"myblog/models"

)

var tmpl *template.Template

func LoadTemplates(){
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActiveLink":isActiveLink,
		"formatDateTime":formatDateTime,
		"postHasTag":postHasTag,
		"noescape":noescape,
		"recentPosts":recentPosts,
		"Excerpt":Excerpt,
		"tags":tags,
		"rankPosts":rankPosts,
		"currentUser":currentUser,
		"pageList":pageList,
		
	})
	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := filepath.Walk("views", fn); err != nil {
		panic(err)
	}
}

func GetTemplates() *template.Template{
	return tmpl
}

func isActiveLink(c *gin.Context,uri string)string{
	if c != nil && c.Request.RequestURI == uri {
		return "active"
	}
	return ""
}

func formatDateTime(t time.Time)string{
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}


func postHasTag(post models.Post,tagTitle string)bool{
	if post.ID == 0 || len(post.Tags) == 0 || len(tagTitle) == 0{
		return false
	}
	for i := range post.Tags{
		if post.Tags[i].Title == tagTitle{
			return true
		}
	}
	return false
}

func noescape(content string)template.HTML{
	return  template.HTML(blackfriday.MarkdownCommon([]byte(content)))
}

func recentPosts(n int)[]models.Post{
	db := models.GetDB()
	var list []models.Post
	if n != 0{
		db.Preload("Tags").Preload("User").Where("published = true").Order("id desc").Limit(n).Find(&list)
	}else{
		db.Preload("Tags").Preload("User").Where("published = true").Order("id desc").Find(&list)
	}
	return list
}

func Excerpt(content string,n int) template.HTML {
	policy := bluemonday.StrictPolicy()
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(content))))
	excerpt := template.HTML(truncate(sanitized, n) + "...")
	return excerpt
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}


func tags() []models.Tag{
	var tags []models.Tag

	models.GetDB().Preload("Posts","published = true").Find(&tags)
	result := make([]models.Tag,0,len(tags))
	for i := range tags {
		if len(tags[i].Posts) > 0 {
			result = append(result,tags[i])
		}
	}
	return result
}

func rankPosts()[]models.Post{
	db := models.GetDB()
	var list []models.Post
	db.Where("published = true").Order("view desc").Limit(4).Find(&list)
	return list
}

func currentUser(c *gin.Context)models.User{

	user := models.User{}
	if user,_ := c.Get("User");user != nil {
		v,ok := user.(*models.User)
		if ok {
			return *v
		}
	}
	return user
}


func pageList(num int) []int {
    ret := make([]int, num+1)
    for i := 1; i < num+1; i++ {
		
    	ret[i] = i 
	}
	ret = ret[1:]
    return ret

}
