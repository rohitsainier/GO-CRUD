package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var Users []User

func main() {
	server := gin.Default()
	server.GET("/", running)
	usersRoutes := server.Group("/users")
	{
		usersRoutes.GET("/list", ListUsers)
		usersRoutes.POST("/create", CreateUser)
		usersRoutes.PUT("/:id", UpdateUser)
		usersRoutes.DELETE("/:id", DeleteUser)
	}
	server.Run()
}

func running(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server is running...",
	})
}
func ListUsers(c *gin.Context) {
	c.JSON(200, Users)
}
func CreateUser(c *gin.Context) {
	var reqBody User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}
	reqBody.ID = uuid.New().String()
	Users = append(Users, reqBody)
	c.JSON(200, gin.H{
		"error":   "false",
		"message": "User added successfully",
		"data":    reqBody,
	})
}
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var reqBody User
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}
	for index, user := range Users {
		if user.ID == id {
			Users[index].ID = id
			Users[index].Firstname = reqBody.Firstname
			Users[index].Lastname = reqBody.Lastname
			c.JSON(200, gin.H{
				"error":   false,
				"message": "User updated successfully",
				"data":    reqBody,
			})
			return
		}
	}
	c.JSON(404, gin.H{
		"error":   true,
		"message": "invalid user id",
	})

}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	for index, user := range Users {
		if user.ID == id {
			Users = append(Users[:index], Users[index+1:]...)
			c.JSON(200, gin.H{
				"error":   false,
				"message": "User deleted successfully",
			})
			return
		}
	}
	c.JSON(404, gin.H{
		"error":   true,
		"message": "invalid user id",
	})

}
