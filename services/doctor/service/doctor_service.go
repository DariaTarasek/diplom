package service

import "github.com/DariaTarasek/diplom/services/doctor/clients"

type DoctorService struct {
	StorageClient *clients.StorageClient
	AuthClient    *clients.AuthClient
}

func NewDoctorService(client *clients.StorageClient, auth *clients.AuthClient) *DoctorService {
	return &DoctorService{
		StorageClient: client,
		AuthClient:    auth,
	}
}
