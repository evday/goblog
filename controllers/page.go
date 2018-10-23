package controllers

// import (
// 	"github.com/gin-gonic/gin"
// )

// const PER_PAGE_COUNT = 10
// const MAX_PAGER_COUNT = 11

// type Pagination struct {
// 	CurrentPage int
// 	TotalCount int
// 	PerPage int
// 	MaxPage int
// 	MaxPageNum int
// 	HalfPage int
// 	BaseUrl string
// 	Context *gin.Context
// }

// func NewPagination(currentPage int,totalCount int,baseUrl string,context *gin.Context) {
// 	pagination := new(Pagination)

// 	if currentPage <= 0{
// 		currentPage = 1
// 	}
// 	pagination.CurrentPage = currentPage // 当前页
// 	pagination.TotalCount = totalCount   // 数据总长度

// 	pagination.PerPage = PER_PAGE_COUNT  //每页显示条数

// 	maxPageNum := totalCount / pagination.PerPage // 页面上应该显示的最大页码
// 	c := totalCount%pagination.PerPage
// 	if c > 0 {
// 		maxPageNum+=1
// 	}
// 	pagination.MaxPageNum = maxPageNum

// 	pagination.MaxPage = MAX_PAGER_COUNT
// 	pagination.HalfPage = (pagination.MaxPage - 1) / 2 // 中间页码

// 	pagination.BaseUrl = baseUrl  //URL 前缀
// 	pagination.Context = context
// }


// func (p *Pagination) Start() int {
// 	return (p.CurrentPage - 1) * p.PerPage
// }

// func (p *Pagination) End() int {
// 	return p.CurrentPage * p.PerPage
// }

// func (p *Pagination) PageHtml() string {
// 	if p.MaxPageNum <= p.MaxPage {
// 		pager_start := 1
// 		pager_end := p.MaxPageNum
// 	}else {
// 		if p.CurrentPage <= p.HalfPage {
// 			pager_start := 1
// 			pager_end := p.MaxPage
// 		}else {
// 			if (p.CurrentPage+p.HalfPage) > p.MaxPageNum {
// 				pager_end := p.MaxPageNum
// 				pager_start := p.MaxPageNum - p.MaxPage + 1
// 			}else {
// 				pager_start := p.CurrentPage - p.HalfPage
// 				pager_end := p.CurrentPage + p.HalfPage
// 			}
//  		}
// 	}

// 	var page_html_list []string

// 	//首页
// 	// page := p.Context.Param("page")
// 	// pagenum,err  := strconv.Atoi(page)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// first_page := fmt.Sprintf(`<li><a href="%s?%s">首页</a></li>`,p.Context.Request.URL(),"sfdsdfljs")
// 	return ""


// }


