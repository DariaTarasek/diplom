package service

import (
	"github.com/DariaTarasek/diplom/services/patient/clients"
)

type PatientService struct {
	StorageClient *clients.StorageClient
	AuthClient    *clients.AuthClient
}

func NewPatientService(client *clients.StorageClient, auth *clients.AuthClient) *PatientService {
	return &PatientService{
		StorageClient: client,
		AuthClient:    auth,
	}
}
