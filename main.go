package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	router.LoadHTMLGlob("static/*")

	s := NewState()

	registerRoutes(s)
	router.Run()
}
