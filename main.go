package main

import (
	"auth-task/routes"
	"auth-task/helpers"
	"fmt"
)

func main() {
	router := routes.InitRoute()
	port := helpers.EnvVar("PORT")
	fmt.Println(port)
	router.Run(port)
}
