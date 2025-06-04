package main

import (
	"github.com/DariaTarasek/diplom/services/statistics/clients"
	grpcserver "github.com/DariaTarasek/diplom/services/statistics/grpc"
	pb "github.com/DariaTarasek/diplom/services/statistics/proto/statistics"
	"github.com/DariaTarasek/diplom/services/statistics/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	storageClient, err := clients.NewStorageClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось создать клиент storage: %s", err)
	}

	statisticsService := service.NewStatisticsService(storageClient)
	lis, err := net.Listen("tcp", ":50056")
	if err != nil {
		log.Fatalf("Не удалось начать слушать: %v", err)
	}

	s := grpc.NewServer()

	server := &grpcserver.Server{
		Service: statisticsService,
	}

	pb.RegisterStatisticsServiceServer(s, server)
	log.Println("Statistics gRPC server started on :50056")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
	log.Println("Сервис статистики запущен.")
}
