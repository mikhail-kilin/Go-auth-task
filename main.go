package main

import (
	"auth-task/routes"
	"auth-task/helpers"
	"fmt"
	"flag"
)

func main() {
	port := flag.String("p", helpers.EnvVar("PORT"), "PORT")
	flag.Parse()

	router := routes.InitRoute()
	fmt.Println(port)
	router.Run(":" + *port)
}
