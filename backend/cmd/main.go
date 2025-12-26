package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!!!",
		})
	})
	if err := router.Run(); err != nil {
		fmt.Printf("server failed: %v", err)
	}
}
