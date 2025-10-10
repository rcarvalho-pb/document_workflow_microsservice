package main

import (
	"context"
	"log"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/adapter/repository/sqlite3"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/api"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/api/proto/userpb"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := sqlite3.ConnectoToDB(":memory:")
	userService := service.NewUserService(db)
	userServer := api.UserGRPCServer{
		UnimplementedUserServiceServer: userpb.UnimplementedUserServiceServer{},
		service:                                               userService,
	}
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10tme.Second)
	defer cancel()
