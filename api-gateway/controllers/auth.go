package controllers

import (
	pb "github.com/Asad2730/product_micro/common/auth"
	"github.com/gin-gonic/gin"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(client pb.AuthServiceClient) *AuthClient {
	return &AuthClient{client: client}
}

func (client *AuthClient) Register(c *gin.Context) {
	var user *pb.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, err.Error())
		return
	}

	res, err := client.client.Register(c, user)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(201, res)
}

func (clinet *AuthClient) Login(c *gin.Context) {
	type loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body *loginDetails
	if err := c.ShouldBindJSON(&body).Error; err != nil {
		c.JSON(500, err)
		return
	}

	req := &pb.LoginRequest{Email: body.Email, Password: body.Password}

	res, err := clinet.client.Login(c, req)
	if err != nil {
		c.JSON(400, res)
		return
	}

	c.JSON(200, res)
}
