package db

import (
	"auth-task/helpers"
    "context"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MgClient struct {
	DB      *mongo.Database
	Client  *mongo.Client
	Context context.Context
}

func ClearDb() {
	GetConnection().DB.Drop(mongoConnection.Context)
}

var mongoConnection *MgClient = nil

func GetConnection() *MgClient {
	ctx := context.Background()
	if mongoConnection == nil {
		client, err := mongo.NewClient(options.Client().ApplyURI(helpers.EnvVar("MONGODB_URI")))
		if err != nil {
			log.Fatal(err)
		}

		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}

		dbName := helpers.EnvVar("DB_NAME")

		mongoConnection = &MgClient{DB: client.Database(dbName), Client: client, Context: ctx}
	}
	return mongoConnection
}
