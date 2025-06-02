package main

import (
	"context"
	"github.com/DariaTarasek/diplom/services/auth/clients"
	grpcserver "github.com/DariaTarasek/diplom/services/auth/grpc"
	pb "github.com/DariaTarasek/diplom/services/auth/proto/auth"
	"github.com/DariaTarasek/diplom/services/auth/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Не удалось получить переменные среды: %w", err)
	}
	storageClient, err := clients.NewStorageClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось создать клиент storage: %s", err)
	}
	redisClient, err := clients.NewRedisClient(ctx)
	if err != nil {
		log.Fatalf("Не удалось создать клиент redis: %s", err.Error())
	}
	smsClient := clients.NewSMSClient()
	//redisClient

	authService := service.NewAuthService(storageClient, redisClient, smsClient)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Не удалось начать слушать: %v", err)
	}

	s := grpc.NewServer()

	server := &grpcserver.Server{
		Service: authService,
	}

	pb.RegisterAuthServiceServer(s, server)

	log.Println("Auth gRPC server started on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
	log.Println("Сервис авторизации запущен.")
}
