
package controllers

import (
	"auth-task/models/entity"
	"auth-task/models/services"
	"auth-task/config/db"
	"log"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (auth *AuthController) Login(c *gin.Context) {
	defer db.CloseConection()
	var loginInfo entity.User

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userService := services.UserService{}
	user, errf := userService.FindUser(&loginInfo)
	if errf != nil {
		c.JSON(401, gin.H{"error": "Not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		c.JSON(402, gin.H{"error": "Email or password is invalid."})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tokens, err := userService.GetTokens(user)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access-token": tokens.AccessToken,
		"access-lifetime": "10 minutes",
		"refresh-token" : tokens.RefreshToken,
		"refresh-token-lifetime" : "7 days",
	})
}

func (auth *AuthController) Refresh(c *gin.Context) {
	defer db.CloseConection()
	tokenString := c.Request.Header.Get("Access-token")
	userService := services.UserService{}
	tokens, err := userService.ReGenerateToken(tokenString)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access-token": tokens.AccessToken,
		"access-lifetime": "10 minutes",
		"refresh-token" : tokens.RefreshToken,
		"refresh-token-lifetime" : "7 days",
	})
}

func (auth *AuthController) DeleteRefreshToken(c *gin.Context) {
	defer db.CloseConection()
	tokenString := c.Request.Header.Get("Access-token")
	refreshService := services.RefreshService{}
	err := refreshService.DeleteSessionByAccessToken(tokenString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"result": "ok"})
	}
}

func (auth *AuthController) DeleteAllRefreshTokennOfUser(c *gin.Context) {
	defer db.CloseConection()
	tokenString := c.Request.Header.Get("Access-token")
	refreshService := services.RefreshService{}
	err := refreshService.DeleteAllSessionsOfUser(tokenString)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"result": "ok"})
	}
}

func (auth *AuthController) Profile(c *gin.Context) {
	defer db.CloseConection()
	user := c.MustGet("user").(*(entity.User))

	c.JSON(200, gin.H{
		"user_name":  user.Name,
		"email":      user.Email,
	})
}

func (auth *AuthController) Signup(c *gin.Context) {
	defer db.CloseConection()

	type signupInfo struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
	}
	var info signupInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(401, gin.H{"error": "Please input all fields"})
		return
	}
	user := entity.User{}
	user.Email = info.Email
	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
		return
	}

	user.Password = string(hash)
	user.Name = info.Name
	userService := services.UserService{}
	err = userService.Create(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"result": "ok"})
	}
	return
}
