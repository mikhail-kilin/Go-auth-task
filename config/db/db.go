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

func (d *MgClient) Close() {
	d.Client.Disconnect(d.Context)
}

var mongoConnection *MgClient = nil

func GetConnection() *MgClient {
	ctx := context.Background()
	if mongoConnection == nil {
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + helpers.EnvVar("DB_CONNECTION_STRING")))
		if err != nil {
			log.Fatal(err)
		}

		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}

		mongoConnection = &MgClient{DB: client.Database(helpers.EnvVar("DB_NAME")), Client: client, Context: ctx}
	}
	return mongoConnection
}
