package db

import (
	"fmt"
	"log"
	"os"

	"github.com/vgo0/gotimepad/models"
	"github.com/glebarez/sqlite" // used to avoid CGO when building
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(path string) {
	var err error

	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	  })

	if err != nil {
		log.Fatalf("Unable to connect to database %s: %s", path, err)
	}		
}

func Migrate() {
	err := DB.AutoMigrate(&models.OTPage{})

	if err != nil {
		log.Fatalf("Unable to migrate database: %s", err)
	}
}

func CheckValidDatabase(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("specified database file does not exist, is %s the correct file?", path)
	}

	return nil
}