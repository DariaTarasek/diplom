package service

import (
	"github.com/DariaTarasek/diplom/services/admin/clients"
)

type AdminService struct {
	StorageClient *clients.StorageClient
}

func NewAdminService(client *clients.StorageClient) *AdminService {
	return &AdminService{
		StorageClient: client,
	}
}
