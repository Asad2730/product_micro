package service

import (
	"context"
	"fmt"
	"log"
	"net"

	auth "github.com/Asad2730/product_micro/common/auth"
	"github.com/Asad2730/product_micro/common/util"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type AuthService struct {
	Db *gorm.DB
	auth.UnimplementedAuthServiceServer
	address string
}

func NewAuthServer(db *gorm.DB, adddress string) *AuthService {
	return &AuthService{Db: db, address: adddress}
}

func (s *AuthService) Login(ctx context.Context, user *auth.LoginRequest) (*auth.AUthResponse, error) {
	var usr *auth.User
	if err := s.Db.First(&usr, "email=", user.Email).Error; err != nil {
		return nil, err
	}

	if err := util.ComparePasswords(usr.Password, user.Password); err != nil {
		return nil, err
	}

	token, err := util.GenerateToken(usr.Id, usr.Email)
	if err != nil {
		return nil, err
	}

	response := &auth.AUthResponse{Id: usr.Id, Username: usr.Username, Email: usr.Email, Token: token}
	return response, nil
}

func (s *AuthService) Register(ctx context.Context, user *auth.User) (*auth.AUthResponse, error) {
	new_user := auth.User{
		Id:       uuid.NewString(),
		Username: user.Username,
		Password: util.HashPassword(user.Password),
		Email:    user.Email,
	}

	if err := s.Db.Create(&new_user).Error; err != nil {
		return nil, err
	}

	token, err := util.GenerateToken(new_user.Id, new_user.Email)
	if err != nil {
		return nil, err
	}

	response := &auth.AUthResponse{Id: new_user.Id, Username: new_user.Username, Email: new_user.Email, Token: token}
	return response, nil
}

func (s *AuthService) Start() error {
	listeners, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	gRPC := grpc.NewServer()
	auth.RegisterAuthServiceServer(gRPC, &AuthService{})
	fmt.Println("gRPC server is running at ", s.address)

	return gRPC.Serve(listeners)
}
