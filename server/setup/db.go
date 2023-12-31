package setup

import (
	"fmt"
	"verve-hrms/internal/schema"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var client *gorm.DB

func GetClient() *gorm.DB {
	if client != nil {
		return client
	}

	pgUser := viper.GetString("POSTGRES_USER")
	if pgUser == "" {
		pgUser = "postgres"
	}
	pgPassword := viper.GetString("POSTGRES_PASSWORD")
	if pgPassword == "" {
		pgPassword = "postgres"
	}
	pgHost := viper.GetString("POSTGRES_HOST")
	if pgHost == "" {
		pgHost = "postgres"
	}
	pgPort := viper.GetString("POSTGRES_PORT")
	if pgPort == "" {
		pgPort = "5432"
	}
	pgDB := viper.GetString("POSTGRES_DB")
	if pgDB == "" {
		pgDB = "verve"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pgHost, pgUser, pgPassword, pgDB, pgPort)

	var err error
	client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true, // ! this is needed to translate postgres errors to gorm errors
	})
	if err != nil {
		panic(err)
	}

	err = client.AutoMigrate(
		&schema.User{},
		&schema.ContactInfo{},
		&schema.EmergencyContact{},
		&schema.Superior{},
		&schema.Subordinate{},
		&schema.Title{},
		&schema.Department{},
		&schema.Location{},
		&schema.JobInfo{},
		&schema.SalaryInfo{},
		&schema.SalaryPayment{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	return client
}
