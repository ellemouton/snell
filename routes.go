package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	articles_db "github.com/ellemouton/snell/articles/db"
)

func registerRoutes(s *State) {
	router.GET("/", s.homeHandler)
	router.GET("/post", s.postArticleFormHandler)
	router.POST("/post", s.saveArticleHandler)
	router.GET("/article/view/:id", s.viewArticleHandler)
	router.GET("/fancy", s.paymentHandler)
	router.GET("/blah", s.checkForAuthHandler)
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

	article, err := articles_db.LookupInfo(s.GetDB(), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	// check auth
	// if found, validate and show page
	// else respond with payment challenge (mac + invoice)

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

func (s *State) paymentHandler(c *gin.Context) {

	str := fmt.Sprintf("LSAT macaroon=\"aslkdjfklsajdkfjsdf\", invoice=\"ln155jgfhgasjklsdfk\"")
	c.Writer.Header().Set("www-authenticate", str)

	c.HTML(
		http.StatusPaymentRequired,
		"payment.html",
		gin.H{
			"title": "payment",
		},
	)
}

var authRegex = regexp.MustCompile("LSAT (.*?):([a-f0-9]{64})")

func (s *State) checkForAuthHandler(c *gin.Context) {
	auth := c.GetHeader("Authorization")

	if auth == "" {
		s.paymentHandler(c)
		return
	}

	if !authRegex.MatchString(auth) {
		s.paymentHandler(c)
		return
	}

	matches := authRegex.FindStringSubmatch(auth)
	if len(matches) != 3 {
		s.paymentHandler(c)
		return
	}

	macBytes, preimageHex := matches[1], matches[2]
	fmt.Fprintln(c.Writer, macBytes, preimageHex)
}
