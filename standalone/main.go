package main

import (
	"log"
	"time"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badorm-example/standalone/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gormDB, err := NewGormDBConnection()
	if err != nil {
		panic(err)
	}

	listOfModels := []any{
		models.Product{},
		models.Company{},
		models.Seller{},
		models.Sale{},
	}

	err = badorm.AutoMigrate(listOfModels, gormDB)
	if err != nil {
		panic(err)
	}

	crudProductService, crudProductRepository := badorm.GetCRUD[models.Product, uuid.UUID](gormDB)

	CreateCRUDObjects(gormDB, crudProductRepository)
	QueryCRUDObjects(crudProductService)
}

func NewGormDBConnection() (*gorm.DB, error) {
	dsn := "user=root password=postgres host=localhost port=26257 sslmode=disable dbname=badaas_db"
	var err error
	retryAmount := 10
	retryTime := 5
	for numberRetry := 0; numberRetry < retryAmount; numberRetry++ {
		database, err := gorm.Open(postgres.Open(dsn))
		if err == nil {
			log.Println("Database connection is active")
			return database, nil
		}

		log.Printf("Database connection failed with error %q\n]", err.Error())
		log.Printf(
			"Retrying database connection %d/%d in %ds\n",
			numberRetry+1, retryAmount, retryTime,
		)
		time.Sleep(time.Duration(retryTime) * time.Second)
	}

	return nil, err
}
