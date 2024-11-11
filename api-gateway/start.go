package apigateway

import (
	"log"

	"github.com/Asad2730/product_micro/api-gateway/controllers"
	"github.com/Asad2730/product_micro/api-gateway/routes"
	auth "github.com/Asad2730/product_micro/common/auth"
	product "github.com/Asad2730/product_micro/common/product"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartApi() {
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	r := gin.Default()
	authClient := auth.NewAuthServiceClient(conn)
	productClient := product.NewProductServiceClient(conn)

	authController := controllers.NewAuthClient(authClient)
	productController := controllers.NewProductClient(productClient)

	routes.AuthRoutes(r, authController)
	routes.ProductRoutes(r, productController)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
