package service

import "github.com/DariaTarasek/diplom/services/auth/clients"

type AuthService struct {
	StorageClient *clients.StorageClient
}

func NewAuthService(client *clients.StorageClient) *AuthService {
	return &AuthService{StorageClient: client}
}
