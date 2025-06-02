package service

import (
	"github.com/DariaTarasek/diplom/services/auth/clients"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	StorageClient *clients.StorageClient
	RedisClient   *redis.Client
	SMSClient     *clients.SMSClient
}

// redis *redis.Client
func NewAuthService(client *clients.StorageClient, redis *redis.Client, sms *clients.SMSClient) *AuthService {
	return &AuthService{
		StorageClient: client,
		RedisClient:   redis,
		SMSClient:     sms,
	}
}
