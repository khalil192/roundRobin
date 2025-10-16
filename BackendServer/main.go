package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	setupRoutes(router.Group(""))

	log.Println("Listening on port", *port)
	if err := router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
