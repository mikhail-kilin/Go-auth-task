package db

import (
	"auth-task/helpers"
	"log"

	"github.com/goonode/mogo"
)

var mongoConnection *mogo.Connection = nil

func GetConnection() *mogo.Connection {
	if mongoConnection == nil {
		connectionString := helpers.EnvVar("DB_CONNECTION_STRING")
		dbName := helpers.EnvVar("DB_NAME")
		config := &mogo.Config{
			ConnectionString: connectionString,
			Database:         dbName,
		}
		mongoConnection, err := mogo.Connect(config)
		if err != nil {
			log.Fatal(err)
		} else {
			return mongoConnection
		}
	}
	return mongoConnection
}
