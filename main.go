package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var Users []User

func main() {
	systray.Run(onReady, onExit)
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
	c.Request.Context().Done()
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

func onReady() {
	model, numberOfCore, frequency, cacheSize := getInfo()
	go func() {
		var result string
		for {
			result = getData()
			systray.SetTitle(result)
		}

	}()
	systray.AddMenuItem(fmt.Sprintf("CPU            : %s", model), "Cpu Model")
	systray.AddMenuItem(fmt.Sprintf("Cores          : %74s", strconv.Itoa(int(numberOfCore))), "Number of core")
	systray.AddMenuItem(fmt.Sprintf("Frequency  : %70s", strconv.Itoa(frequency)), "Frequency CPU")
	systray.AddMenuItem(fmt.Sprintf("CPU Cache : %71s", strconv.Itoa(int(cacheSize))), "CPU Cache")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}

		}

	}()
}

func onExit() {

}

func getMemoryUsage() int {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}
	return int(math.Ceil(memory.UsedPercent))
}

func getCpuUsage() int {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatal(err)
	}
	return int(math.Ceil(percent[0]))
}

func getData() string {
	cpuData := "Cpu: " + strconv.Itoa(getCpuUsage()) + "% "
	memoryData := "Ram: " + strconv.Itoa(getMemoryUsage()) + "% "
	return cpuData + memoryData
}

func getInfo() (string, int32, int, int32) {
	info, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}
	return info[0].ModelName, info[0].Cores, int(info[0].Mhz), info[0].CacheSize
}
