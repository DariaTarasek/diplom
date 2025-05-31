package main

import (
	"github.com/DariaTarasek/diplom/services/doctor/clients"
	grpcserver "github.com/DariaTarasek/diplom/services/doctor/grpc"
	pb "github.com/DariaTarasek/diplom/services/doctor/proto/doctor"
	"github.com/DariaTarasek/diplom/services/doctor/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	storageClient, err := clients.NewStorageClient("localhost:50051")
	if err != nil {
		log.Fatalf("Не удалось создать клиент storage: %s", err)
	}
	authClient, err := clients.NewAuthClient("localhost:50052")
	if err != nil {
		log.Fatalf("Не удалось создать клиент auth: %s", err)
	}
	doctorService := service.NewDoctorService(storageClient, authClient)
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Не удалось начать слушать: %v", err)
	}

	s := grpc.NewServer()

	server := &grpcserver.Server{
		Service: doctorService,
	}

	pb.RegisterDoctorServiceServer(s, server)

	log.Println("Doctor gRPC server started on :50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
	log.Println("Сервис авторизации запущен.")
}
