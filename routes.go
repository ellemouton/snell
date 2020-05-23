package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	articles_db "github.com/ellemouton/snell/articles/db"
)

func registerRoutes(s *State) {
	router.GET("/", s.homeHandler)
	router.GET("/post", s.postArticleFormHandler)
	router.POST("/post", s.saveArticleHandler)
	router.GET("/article/view/:id", s.viewArticleHandler)
}

func (s *State) homeHandler(c *gin.Context) {

	articles, err := articles_db.ListAllInfo(s.GetDB())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":   "Home",
			"payload": articles,
		},
	)
}

func (s *State) postArticleFormHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"post.html",
		gin.H{
			"title": "Post Article",
		},
	)
}

func (s *State) saveArticleHandler(c *gin.Context) {
	title := c.PostForm("title")
	abstract := c.PostForm("abstract")
	content := c.PostForm("content")
	fmt.Println(c.PostForm("price"))
	price, err := strconv.ParseInt(c.PostForm("price"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	// TODO(elle): Create ops package for articles rather than calling db layer
	_, err = articles_db.Create(s.GetDB(), title, abstract, price, content)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	s.homeHandler(c)
}

func (s *State) viewArticleHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	fmt.Println(id)

	article, err := articles_db.LookupInfo(s.GetDB(), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	content, err := articles_db.LookupContent(s.GetDB(), article.ContentID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.HTML(
		http.StatusOK,
		"article.html",
		gin.H{
			"title":   "View",
			"payload": content,
		},
	)

}
