package authservice

import (
	"fmt"

	"github.com/Asad2730/product_micro/auth-service/service"
	user "github.com/Asad2730/product_micro/common/auth"
	"github.com/Asad2730/product_micro/database"
)

func StartMicroService() {
	db, err := database.NewPostgresDB()
	if err != nil {
		fmt.Println("Failed Connect to Database", err.Error())
	}
	db.AutoMigrate(&user.User{})
	service := service.NewAuthServer(db, ":8000")

	if err := service.Start(); err != nil {
		fmt.Println("Failed to start Auth service", err.Error())
	}
}
