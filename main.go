package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

func main() {
	flag.Parse()

	router = gin.Default()
	router.LoadHTMLGlob("static/*.html")
	router.Static("/css", "./static/css")
	router.Static("/js", "./static/js")

	s, err := NewState()
	if err != nil {
		log.Fatal(fmt.Sprintf("problem initializing the state: %s", err.Error()))
	}
	defer s.cleanup()

	registerRoutes(s)
	router.Run(":8000")
}
