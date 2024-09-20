package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	PostgresURL string
	MongoDbURI  string
	GormDB      *gorm.DB
	MongoDB     *mongo.Database
	MongoClient *mongo.Client
	Database    string
)

func MainConfiguration(configFile string) error {
	viper.SetConfigFile(configFile)

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config.yaml file: %w", err)
	}

	PostgresURL = viper.GetString("postgres_url")
	MongoDbURI = viper.GetString("mongodb_uri")

	log.Println("Main Configuration Done!!")

	return nil

}
