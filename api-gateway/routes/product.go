package routes

import (
	"github.com/Asad2730/product_micro/api-gateway/controllers"
	"github.com/Asad2730/product_micro/common/util"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine, h *controllers.ProductClient) {
	product := r.Group("/products")
	{
		product.Use(util.VerifyToken())
		product.POST("/", h.CreateProduct)
		product.GET("/", h.GetProduct)
		product.GET("/", h.ListProducts)
		product.PUT("/:id", h.UpdateProduct)
		product.DELETE("/:id", h.DeleteProduct)
	}

	category := r.Group("/category")
	{
		category.Use(util.VerifyToken())
		category.POST("/", h.CreateCategory)
		category.GET("/", h.GetProduct)
		category.GET("/", h.ListCategories)
		category.PUT("/:id", h.UpdateCategory)
		category.DELETE("/:id", h.DeleteCategory)
	}
}
