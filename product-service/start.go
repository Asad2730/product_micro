package productservice

import (
	"fmt"
	"log"
	"os"

	product "github.com/Asad2730/product_micro/common/product"
	"github.com/Asad2730/product_micro/database"
	"github.com/Asad2730/product_micro/product-service/cached"
	"github.com/Asad2730/product_micro/product-service/service"
	"github.com/joho/godotenv"
)

func StartMicroService() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is not set in the .env file")
	}

	db, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Failed Connect to Database", err.Error())
	}
	db.AutoMigrate(&product.Product{}, &product.Category{})
	redis := cached.NewRedis(redisURL)
	service := service.NewProductServer(db, redis, ":8000")

	if err := service.Start(); err != nil {
		fmt.Println("Failed to start Product service", err.Error())
	}
}
