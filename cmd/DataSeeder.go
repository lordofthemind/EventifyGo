package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/lordofthemind/EventifyGo/configs"
	"github.com/lordofthemind/EventifyGo/internals/types"
	"github.com/lordofthemind/mygopher/gopherlogger"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gopherpostgres"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

// Seeder interface for flexibility between Postgres and MongoDB
type Seeder interface {
	SeedSuperUsers(count int)
	SeedEvents(count int)
}

// PostgresSeeder implementation
type PostgresSeeder struct {
	db *gorm.DB
}

func (ps *PostgresSeeder) SeedSuperUsers(count int) {
	for i := 0; i < count; i++ {
		superUser := types.SuperUserType{
			ID:               uuid.New(),
			Role:             "guest",
			Email:            faker.Email(),
			FullName:         faker.Name(),
			Username:         faker.Username(),
			HashedPassword:   faker.Password(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Is2FAEnabled:     false,
			PermissionGroups: []string{"group1", "group2"},
		}
		ps.db.Create(&superUser)
	}
	fmt.Println("Postgres: Seeded SuperUsers")
}

func (ps *PostgresSeeder) SeedEvents(count int) {
	for i := 0; i < count; i++ {
		event := types.EventType{
			EventID:     uuid.New(),
			Name:        faker.Word(),
			Description: faker.Sentence(),
			Date:        time.Now().AddDate(0, 0, rand.Intn(365)),
			Location:    faker.Word(),
			Capacity:    rand.Intn(451) + 50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			OrganizerID: uuid.New(),
			Attendees:   []uuid.UUID{uuid.New(), uuid.New()}, // Correct syntax for UUID array
		}
		ps.db.Create(&event)
	}
	fmt.Println("Postgres: Seeded Events")
}

// MongoSeeder implementation
type MongoSeeder struct {
	db *mongo.Database
}

func (ms *MongoSeeder) SeedSuperUsers(count int) {
	collection := ms.db.Collection("superusers")
	for i := 0; i < count; i++ {
		superUser := types.SuperUserType{
			ID:               uuid.New(),
			Role:             "guest",
			Email:            faker.Email(),
			FullName:         faker.Name(),
			Username:         faker.Username(),
			HashedPassword:   faker.Password(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Is2FAEnabled:     false,
			PermissionGroups: []string{"group1", "group2"},
		}
		collection.InsertOne(context.TODO(), superUser)
	}
	fmt.Println("MongoDB: Seeded SuperUsers")
}

func (ms *MongoSeeder) SeedEvents(count int) {
	collection := ms.db.Collection("events")
	for i := 0; i < count; i++ {
		event := types.EventType{
			EventID:     uuid.New(),
			Name:        faker.Word(),
			Description: faker.Sentence(),
			Date:        time.Now().AddDate(0, 0, rand.Intn(365)),
			Location:    faker.Word(),
			Capacity:    rand.Intn(451) + 50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			OrganizerID: uuid.New(),
			Attendees:   []uuid.UUID{uuid.New(), uuid.New()},
		}
		collection.InsertOne(context.TODO(), event)
	}
	fmt.Println("MongoDB: Seeded Events")
}

// SeederFactory for creating appropriate Seeder instance
func SeederFactory(dbType string, db interface{}) Seeder {
	switch dbType {
	case "postgres":
		return &PostgresSeeder{db: db.(*gorm.DB)}
	case "mongodb":
		return &MongoSeeder{db: db.(*mongo.Database)}
	}
	return nil
}

func main() {

	logFile, err := gopherlogger.SetUpLoggerFile("Seeder.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()

	err = configs.SeederConfiguration("config.yaml")
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

	mongoDB := gophermongo.GetDatabase(mongoClient, "superuser")

	// Choose the database to seed
	seeder := SeederFactory("mongodb", mongoDB) // Change to "mongodb" for MongoDB seeding

	// Seed SuperUsers and Events
	seeder.SeedSuperUsers(10)
	seeder.SeedEvents(5)
}
