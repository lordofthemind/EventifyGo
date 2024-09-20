package cmd

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/EventifyGo/configs"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gopherpostgres"
	"github.com/lordofthemind/mygopher/mygopherlogger"
)

func GinServer() {
	logFile, err := mygopherlogger.SetUpLoggerFile("ginServer.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()

	err = configs.MainConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Postgres connection
	ctx := context.Background()
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
		log.Fatalf("failed to migrate Postgres database: %v", err)
	}

	// MongoDB connection
	mongoClient, err := gophermongo.ConnectToMongoDB(ctx, configs.MongoDbURI, 10*time.Second, 3)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	_ = gophermongo.GetDatabase(mongoClient, "superuser")
}
