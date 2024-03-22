package main

import (
	"os"

	"github.com/1chickin/authen-jwt-redis/config"
	"github.com/1chickin/authen-jwt-redis/controller"
	"github.com/1chickin/authen-jwt-redis/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
	// config.MigrateDB()
}

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping pong",
		})
	})

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)
	r.GET("/validate-token", middleware.RequireAuth, controller.ValidateToken)
	r.Run(":" + os.Getenv("PORT"))

}
