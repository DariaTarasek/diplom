package service

import (
	"github.com/DariaTarasek/diplom/services/patient/clients"
)

type PatientService struct {
	StorageClient *clients.StorageClient
}

func NewPatientService(client *clients.StorageClient) *PatientService {
	return &PatientService{
		StorageClient: client,
	}
}
