package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/Asad2730/product_micro/common/product"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type ProductServiceServer struct {
	Db      *gorm.DB
	cache   *redis.Client
	address string
	pb.UnimplementedProductServiceServer
}

func NewProductServer(db *gorm.DB, cache *redis.Client, address string) *ProductServiceServer {
	return &ProductServiceServer{Db: db, cache: cache, address: address}
}

func (s *ProductServiceServer) CreateCategory(ctx context.Context, req *pb.Category) (*pb.CategoryResponse, error) {

	new_category := pb.Category{
		Id:   uuid.NewString(),
		Name: req.Name,
	}

	if err := s.Db.Create(&new_category).Error; err != nil {
		return nil, err
	}

	response := &pb.CategoryResponse{Category: &new_category}
	if err := s.cache.Del(ctx, "category").Err(); err != nil {
		return nil, err
	}
	return response, nil

}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.Product) (*pb.ProductResponse, error) {
	new_product := pb.Product{
		Id:         uuid.NewString(),
		Name:       req.Name,
		Price:      req.Price,
		CategoryId: req.CategoryId,
	}

	if err := s.Db.Create(&new_product).Error; err != nil {
		return nil, err
	}

	response := &pb.ProductResponse{Product: &new_product}
	if err := s.cache.Del(ctx, "product").Err(); err != nil {
		return nil, err
	}
	return response, nil

}
func (s *ProductServiceServer) DeleteCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {

	var category *pb.Category

	if err := s.Db.First(&category, pb.CategoryRequest{Id: req.Id}).Error; err != nil {
		return nil, err
	}

	if err := s.Db.Delete(&pb.Category{}, req.Id).Error; err != nil {
		return nil, err
	}
	s.cache.Del(ctx, "category/"+req.Id)

	response := &pb.CategoryResponse{Category: category}
	return response, nil
}
func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {

	var product *pb.Product
	if err := s.Db.First(&product, req.Id).Error; err != nil {
		return nil, err
	}

	if err := s.Db.Delete(&pb.Product{}, req.Id).Error; err != nil {
		return nil, err
	}
	s.cache.Del(ctx, "product/"+req.Id)

	response := &pb.ProductResponse{Product: product}
	return response, nil

}

func (s *ProductServiceServer) GetCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	cacheKey := "category/" + req.Id

	cachedResponse, err := s.cache.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedResponse != nil {
		var response pb.CategoryResponse
		if err := proto.Unmarshal(cachedResponse, &response); err == nil {
			return &response, nil
		}
	}

	var category *pb.Category
	if err := s.Db.First(&category, req.Id).Error; err != nil {
		return nil, err
	}

	response := &pb.CategoryResponse{Category: category}
	encodedResponse, err := proto.Marshal(response)
	if err == nil {
		s.cache.Set(ctx, cacheKey, encodedResponse, 15*time.Minute)
	}

	return response, nil
}

func (s *ProductServiceServer) GetProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	cacheKey := "product/" + req.Id
	cachedResponse, err := s.cache.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedResponse != nil {
		var response pb.ProductResponse
		if err := proto.Unmarshal(cachedResponse, &response); err == nil {
			return &response, nil
		}
	}

	var product *pb.Product

	if err := s.Db.First(&product, req.Id).Error; err != nil {
		return nil, err
	}

	response := &pb.ProductResponse{Product: product}
	encodedResponse, err := proto.Marshal(response)
	if err == nil {
		s.cache.Set(ctx, cacheKey, encodedResponse, 15*time.Minute)
	}

	return response, nil

}

func (s *ProductServiceServer) ListCategories(ctx context.Context, req *pb.CategoryListRequest) (*pb.CategoryListResponse, error) {
	cacheKey := "categories_page_" + strconv.Itoa(int(req.Page)) + "_size_" + strconv.Itoa(int(req.PageSize))
	cachedResponse, err := s.cache.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedResponse != nil {
		var response pb.CategoryListResponse
		if err := proto.Unmarshal(cachedResponse, &response); err == nil {
			return &response, nil
		}
	}

	var categories []*pb.Category

	page := req.Page
	pageSize := req.PageSize
	offset := (page - 1) * pageSize

	if err := s.Db.Limit(int(pageSize)).Offset(int(offset)).Find(&categories).Error; err != nil {
		return nil, err
	}

	res := &pb.CategoryListResponse{Categories: categories}

	encodedResponse, err := proto.Marshal(res)
	if err == nil {
		s.cache.Set(ctx, cacheKey, encodedResponse, 15*time.Minute)
	}

	return res, nil
}

func (s *ProductServiceServer) ListProducts(ctx context.Context, req *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	cacheKey := "products_page_" + strconv.Itoa(int(req.Page)) + "_size_" + strconv.Itoa(int(req.PageSize))
	cachedResponse, err := s.cache.Get(ctx, cacheKey).Bytes()
	if err == nil && cachedResponse != nil {
		var response pb.ProductListResponse
		if err := proto.Unmarshal(cachedResponse, &response); err == nil {
			return &response, nil
		}
	}

	var products []*pb.Product
	var categories []*pb.Category

	page := req.Page
	pageSize := req.PageSize
	offset := (page - 1) * pageSize

	if err := s.Db.Limit(int(pageSize)).Offset(int(offset)).Find(&products).Error; err != nil {
		return nil, err
	}

	for _, product := range products {
		var categorie *pb.Category
		if err := s.Db.First(&categorie, product.CategoryId).Error; err != nil {
			return nil, err
		}

		categories = append(categories, categorie)
	}
	response := &pb.ProductListResponse{Products: products, Categories: categories}

	return response, nil
}

func (s *ProductServiceServer) UpdateCategory(ctx context.Context, req *pb.Category) (*pb.CategoryResponse, error) {
	var category pb.Category

	if err := s.Db.First(&category, req.Id).Error; err != nil {
		return nil, err
	}

	category.Name = req.Name

	if err := s.Db.Save(&category).Error; err != nil {
		return nil, err
	}

	s.cache.Del(ctx, "category/"+req.Id)

	return &pb.CategoryResponse{Category: &category}, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.ProductResponse, error) {
	var product pb.Product

	if err := s.Db.First(&product, req.Id).Error; err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Price = req.Price
	product.CategoryId = req.CategoryId

	if err := s.Db.Save(&product).Error; err != nil {
		return nil, err
	}

	s.cache.Del(ctx, "product/"+req.Id)

	return &pb.ProductResponse{Product: &product}, nil
}

func (s *ProductServiceServer) Start() error {
	listeners, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	gRPC := grpc.NewServer()
	pb.RegisterProductServiceServer(gRPC, &ProductServiceServer{})
	fmt.Println("gRPC server is running at ", s.address)

	return gRPC.Serve(listeners)
}
