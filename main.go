package main

import (
	apigateway "github.com/Asad2730/product_micro/api-gateway"
	authservice "github.com/Asad2730/product_micro/auth-service"
	productservice "github.com/Asad2730/product_micro/product-service"
)

func main() {
	authservice.StartMicroService()
	productservice.StartMicroService()
	apigateway.StartApi()
}
