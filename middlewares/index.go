package middlewares

import (
	"fmt"

	"auth-task/models/entity"
	"auth-task/models/services"
	"auth-task/helpers"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
			token_time, errt := time.Parse(time.RFC3339, claims["created_at"].(string))
			if errt != nil {
				c.AbortWithStatusJSON(402, gin.H{
					"error": "Invalid Token",
				})
				return
			}

			token_time = token_time.Add(10 * time.Minute)
			current_time := time.Now()
			if current_time.After(token_time) {
				c.AbortWithStatusJSON(401, gin.H{
					"error": "Token is outdated",
				})
				return
			}

			email := claims["email"].(string)
			fmt.Println("email is ", email)
			userservice := services.Userservice{}
			user, err := userservice.FindUser(&entity.User {Email: email})
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
