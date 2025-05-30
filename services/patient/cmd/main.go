package main

import (
	"github.com/DariaTarasek/diplom/services/patient/clients"
	grpcserver "github.com/DariaTarasek/diplom/services/patient/grpc"
	pb "github.com/DariaTarasek/diplom/services/patient/proto/patient"
	"github.com/DariaTarasek/diplom/services/patient/service"
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
		log.Fatalf("Не удалось создать клиент auth: %s", err.Error())
	}

	patientService := service.NewPatientService(storageClient, authClient)
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Не удалось начать слушать: %v", err)
	}

	s := grpc.NewServer()

	server := &grpcserver.Server{
		Service: patientService,
	}

	pb.RegisterPatientServiceServer(s, server)

	log.Println("Patient gRPC server started on :50054")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
	log.Println("Сервис авторизации запущен.")
}
