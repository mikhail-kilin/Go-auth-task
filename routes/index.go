package routes

import (
	"auth-task/controllers"

	//"auth-task/middlewares"

	"github.com/gin-gonic/gin"
)

func setAuthRoute(router *gin.Engine) {
	authController := new(controllers.AuthController)
	//router.POST("/login", authController.Login)
	router.POST("/signup", authController.Signup)

	//authGroup := router.Group("/")
	//authGroup.Use(middlewares.Authentication())
	//authGroup.GET("/profile", authController.Profile)

}

func InitRoute() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	setAuthRoute(router)
	return router
}