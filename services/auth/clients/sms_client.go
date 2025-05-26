package clients

import (
	smsaero_golang "github.com/smsaero/smsaero_golang/smsaero"
	"os"
)

type SMSClient struct {
	Client *smsaero_golang.Client
}

func NewSMSClient() *SMSClient {
	username := os.Getenv("SMSAERO_USERNAME")
	apiKey := os.Getenv("SMSAERO_APIKEY")
	client := smsaero_golang.NewSmsAeroClient(username, apiKey, smsaero_golang.WithPhoneValidation(false))
	return &SMSClient{Client: client}
}
