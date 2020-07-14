package middlewares

import (
	"fmt"

	"auth-task/models/entity"
	"auth-task/models/services"
	"auth-task/helpers"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Access-token")
		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Access-token header is missing",
			})
			return
		}
		fmt.Println("tokenString is ", tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secretKey := helpers.EnvVar("SECRET")
			return []byte(secretKey), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			tokenTime, errt := time.Parse(time.RFC3339, claims["created_at"].(string))
			if errt != nil {
				c.AbortWithStatusJSON(402, gin.H{
					"error": "Invalid Token",
				})
				return
			}

			tokenTime = tokenTime.Add(10 * time.Minute)
			currentTime := time.Now()
			if currentTime.After(tokenTime) {
				c.AbortWithStatusJSON(401, gin.H{
					"error": "Token is outdated",
				})
				return
			}

			email := claims["email"].(string)
			fmt.Println("email is ", email)
			userService := services.UserService{}
			user, err := userService.FindUser(&entity.User {Email: email})
			if err != nil {
				c.AbortWithStatusJSON(402, gin.H{
					"error": "User not found",
				})
				return
			}
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Token is not valid",
			})
		}
	}
}

func AuthRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Access-token")
		refreshTokenString := c.Request.Header.Get("Refresh-token")

		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Access-token header is missing",
			})
			return
		}

		if len(refreshTokenString) == 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Refresh-token header is missing",
			})
			return
		}

		secretKey := helpers.EnvVar("SECRET")
		fmt.Println("tokenString is ", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email := claims["email"].(string)
			fmt.Println("email is ", email)
			userService := services.UserService{}
			user, err := userService.FindUser(&entity.User {Email: email})
			if err != nil {
				c.AbortWithStatusJSON(402, gin.H{
					"error": "User not found",
				})
				return
			}
			c.Set("user", user)

			sessionId := claims["session_id"].(string)

			refreshService := services.RefreshService{}
			session, errs := refreshService.FindSession(sessionId)

			if errs != nil {
				c.AbortWithStatusJSON(401, gin.H{
					"error": "Refresh Token is invalid",
				})
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(session.Refresh), []byte(refreshTokenString))
			if err != nil {
				c.AbortWithStatusJSON(402, gin.H{"error": "Refresh Token is invalid."})
				return
			}

			refreshTokenTime := session.CreatedAt.Add(7 * 24 * time.Hour)
			currentTime := time.Now()
			if currentTime.After(refreshTokenTime) {
				refreshService.DeleteSession(sessionId)
				c.AbortWithStatusJSON(401, gin.H{
					"error": "Refresh Token is outdated",
				})
				return
			}
			c.Next()
		} else {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Token is not valid",
			})
		}
	}
}
