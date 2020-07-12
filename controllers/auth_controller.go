
package controllers

import (
	"auth-task/models/entity"
	"auth-task/models/services"
	"log"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (auth *AuthController) Login(c *gin.Context) {

	var loginInfo entity.User

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userservice := services.Userservice{}
	user, errf := userservice.FindUser(&loginInfo)
	if errf != nil {
		c.JSON(401, gin.H{"error": "Not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		c.JSON(402, gin.H{"error": "Email or password is invalid."})
		return
	}

	token, err := user.GetJwtToken()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	refresh_service := services.RefreshService{}
	refrsh_token, err := refresh_service.Generate(user, token)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access-token": token,
		"access-lifetime": "10 minutes",
		"refresh-token" : refrsh_token,
		"refresh-token-lifetime" : "7 days",
	})
}

func (auth *AuthController) Profile(c *gin.Context) {
	user := c.MustGet("user").(*(entity.User))

	c.JSON(200, gin.H{
		"user_name":  user.Name,
		"email":      user.Email,
	})
}

func (auth *AuthController) Signup(c *gin.Context) {

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
	userservice := services.Userservice{}
	err = userservice.Create(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"result": "ok"})
	}
	return
}
