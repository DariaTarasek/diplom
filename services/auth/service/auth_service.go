package service

import (
	"github.com/DariaTarasek/diplom/services/auth/clients"
)

type AuthService struct {
	StorageClient *clients.StorageClient
	//RedisClient   *redis.Client
	SMSClient *clients.SMSClient
}

// redis *redis.Client
func NewAuthService(client *clients.StorageClient, sms *clients.SMSClient) *AuthService {
	return &AuthService{
		StorageClient: client,
		//RedisClient:   redis,
		SMSClient: sms,
	}
}
