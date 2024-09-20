package cmd

import (
	"log"

	"github.com/lordofthemind/EventifyGo/configs"
	"github.com/lordofthemind/EventifyGo/internals/handlers"
	"github.com/lordofthemind/EventifyGo/internals/initializers"
	"github.com/lordofthemind/EventifyGo/internals/repositories"
	"github.com/lordofthemind/EventifyGo/internals/repositories/mongodb"
	"github.com/lordofthemind/EventifyGo/internals/services"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/mygopherlogger"
	"github.com/revel/revel"
)

func RevelServer() {
	// Set up logger
	logFile, err := mygopherlogger.SetUpLoggerFile("revelServer.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()

	// Load configuration
	err = configs.MainConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Initialize database (Postgres or MongoDB)
	initializers.DatabaseInitializer()

	// Setup repository and service based on the selected database
	var superUserRepository repositories.SuperUserRepositoryInterface

	switch configs.Database {
	case "postgres":
		if configs.GormDB == nil {
			log.Fatalf("Postgres connection was not initialized")
		}
		// superUserRepository = postgres.NewPostgresSuperUserRepository(configs.GormDB) // Example
	case "mongodb":
		if configs.MongoClient == nil {
			log.Fatalf("MongoDB client was not initialized")
		}
		superUserDB := gophermongo.GetDatabase(configs.MongoClient, "superuser")
		superUserRepository = mongodb.NewMongoSuperUserRepository(superUserDB)
	default:
		log.Fatalf("Invalid database configuration")
	}

	// Initialize service
	superUserService := services.NewSuperUserService(superUserRepository)

	// Register controller with the service
	superUserController := new(handlers.SuperUserRevelController)
	superUserController.Init(superUserService)

	// Start the Revel server
	// Note: Revel expects to be run via its own command. The following is a simplified example.
	revel.Run(8080)
}
