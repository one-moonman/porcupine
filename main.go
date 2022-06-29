package main

import (
	"bug-free-octo-broccoli/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/register", controllers.Register())
	router.POST("/login", controllers.Login())

	// forAccessTokenVerification := router.Group("/")
	// forAccessTokenVerification.Use()
	// {
	// 	forAccessTokenVerification.POST("/logout")
	// 	forAccessTokenVerification.DELETE("/delete")
	// }
	router.Run()
}
