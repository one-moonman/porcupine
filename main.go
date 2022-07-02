package main

import "bug-free-octo-broccoli/api"

func main() {
	// router := gin.Default()

	// router.POST("/register", controllers.Register())
	// router.POST("/login", controllers.Login())

	// forAccessTokenVerification := router.Group("/")
	// forAccessTokenVerification.Use(middlewares.VerifyAccessToken())
	// {
	// 	forAccessTokenVerification.GET("/me", controllers.Me())
	// 	forAccessTokenVerification.POST("/logout", controllers.Logout())
	// 	forAccessTokenVerification.DELETE("/delete")
	// }

	// router.Run()
	server := new(api.Server)
	server.Run()
}
