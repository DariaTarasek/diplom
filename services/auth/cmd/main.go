package main

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/clients"
	grpcserver "github.com/DariaTarasek/diplom/services/auth/grpc"
	pb "github.com/DariaTarasek/diplom/services/auth/proto/auth"
	"github.com/DariaTarasek/diplom/services/auth/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
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
	authService := service.NewAuthService(storageClient, redisClient, smsClient)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("не удалось начать слушать: %v", err)
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

func WaitForTCP(address string, retryInterval time.Duration) {
	fmt.Printf("⏳ Waiting for %s...\n", address)
	for {
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err == nil {
			_ = conn.Close()
			fmt.Printf("✅ %s is available\n", address)
			return
		}

		// Проверка, не завершён ли процесс
		select {
		case <-time.After(retryInterval):
			// продолжаем ждать
		default:
			// можно вставить логику выхода, если нужно
		}

		fmt.Printf("❌ Still waiting for %s (error: %v)\n", address, err)
	}
}
