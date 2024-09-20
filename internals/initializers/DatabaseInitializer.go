package initializers

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/EventifyGo/configs"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gopherpostgres"
)

func DatabaseInitializer() {
	ctx := context.Background()

	if configs.Database == "postgres" {
		// Initialize PostgreSQL
		gormDB, err := gopherpostgres.ConnectToPostgresGORM(ctx, configs.PostgresURL, 10*time.Second, 3)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL using GORM: %v", err)
		}

		err = gopherpostgres.CheckAndEnableUUIDExtension(gormDB)
		if err != nil {
			log.Fatalf("Failed to confirm UUID extension: %v", err)
		}

		// Auto migrate for GORM (Postgres)
		if err := gormDB.AutoMigrate(&types.SuperUserType{}, &types.EventType{}); err != nil {
			log.Fatalf("Failed to migrate Postgres database: %v", err)
		}

		// Set global GormDB
		configs.GormDB = gormDB
	}

	if configs.Database == "mongodb" {
		// Initialize MongoDB
		mongoClient, err := gophermongo.ConnectToMongoDB(ctx, configs.MongoDbURI, 10*time.Second, 3)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Set global MongoDB and MongoClient
		configs.MongoClient = mongoClient
	}
}
