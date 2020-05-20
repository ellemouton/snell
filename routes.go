package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(s *State) {
	router.GET("/", s.homeHandler)
	//router.GET("/post", s.postArticleHandler)
}

func (s *State) homeHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home.html",
		gin.H{
			"title": "Tadda!",
		},
	)
}
