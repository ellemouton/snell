package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	router.LoadHTMLGlob("static/*")

	s, err := NewState()
	if err != nil {
		log.Fatal(fmt.Sprintf("problem initializing the state: %s", err.Error()))
	}

	registerRoutes(s)
	router.Run()
}
