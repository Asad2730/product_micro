package controllers

import (
	pb "github.com/Asad2730/product_micro/common/product"
	"github.com/gin-gonic/gin"
)

type ProductClient struct {
	client pb.ProductServiceClient
}

func NewProductClient(client pb.ProductServiceClient) *ProductClient {
	return &ProductClient{client: client}
}

func (clinet *ProductClient) CreateProduct(c *gin.Context) {
	var product pb.Product
	if err := c.ShouldBindJSON(&product).Error; err != nil {
		c.JSON(500, err)
		return
	}

	res, err := clinet.client.CreateProduct(c, &product)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, res)
}

func (client *ProductClient) GetProduct(c *gin.Context) {
	id := c.Param("id")
	res, err := client.client.GetProduct(c, &pb.ProductRequest{Id: id})
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, res)
}

func (client *ProductClient) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var body pb.Product
	if err := c.ShouldBindJSON(&body).Error; err != nil {
		c.JSON(500, err)
		return
	}

	res, err := client.client.UpdateProduct(c, &pb.Product{Id: id, Name: body.Name, Price: body.Price, CategoryId: body.CategoryId})
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, res)
}

func (client *ProductClient) ListProducts(c *gin.Context) {
	res, err := client.client.ListProducts(c, &pb.ProductListRequest{})
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, res)
}

func (client *ProductClient) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := client.client.DeleteProduct(c, &pb.ProductRequest{Id: id})
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, "Product deleted successfully")
}

func (client *ProductClient) CreateCategory(c *gin.Context) {
	var category pb.Category
	if err := c.ShouldBindJSON(&category).Error; err != nil {
		c.JSON(500, err)
		return
	}

	res, err := client.client.CreateCategory(c, &category)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(201, res)
}

func (client *ProductClient) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var body pb.Category
	if err := c.ShouldBindJSON(&body).Error; err != nil {
		c.JSON(500, err)
		return
	}

	res, err := client.client.UpdateCategory(c, &pb.Category{Id: id, Name: body.Name})
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, res)
}

func (client *ProductClient) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	_, err := client.client.DeleteCategory(c, &pb.CategoryRequest{Id: id})
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, "Category deleted successfully")
}

func (client *ProductClient) ListCategories(c *gin.Context) {
	res, err := client.client.ListCategories(c, &pb.CategoryListRequest{})
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, res)
}
