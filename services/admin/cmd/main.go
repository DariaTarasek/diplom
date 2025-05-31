package main

import (
	"github.com/DariaTarasek/diplom/services/admin/clients"
	grpcserver "github.com/DariaTarasek/diplom/services/admin/grpc"
	pb "github.com/DariaTarasek/diplom/services/admin/proto/admin"
	"github.com/DariaTarasek/diplom/services/admin/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	storageClient, err := clients.NewStorageClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось создать клиент storage: %s", err)
	}

	adminService := service.NewAdminService(storageClient)
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Не удалось начать слушать: %v", err)
	}

	s := grpc.NewServer()

	server := &grpcserver.Server{
		Service: adminService,
	}

	pb.RegisterAdminServiceServer(s, server)
	log.Println("Admin gRPC server started on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
	log.Println("Сервис администратора запущен.")
}
