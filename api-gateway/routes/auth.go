package routes

import (
	"github.com/Asad2730/product_micro/api-gateway/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, h *controllers.AuthClient) {
	user := r.Group("/auth")
	{
		user.POST("/login", h.Login)
		user.POST("/signup", h.Register)
	}
}
