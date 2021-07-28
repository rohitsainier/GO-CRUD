package main

import "github.com/gin-gonic/gin"

type Person struct {
	ID        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var people []Person

func main() {
	server := gin.Default()
	server.GET("/", running)
	server.Run()
}

func running(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server is running...",
	})
}
