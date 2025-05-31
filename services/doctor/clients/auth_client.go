package clients

import (
	authpb "github.com/DariaTarasek/diplom/services/doctor/proto/auth"
	"google.golang.org/grpc"
)

type AuthClient struct {
	Conn   *grpc.ClientConn
	Client authpb.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure()) // или grpc.WithTransportCredentials
	if err != nil {
		return nil, err
	}

	client := authpb.NewAuthServiceClient(conn)
	return &AuthClient{
		Conn:   conn,
		Client: client,
	}, nil
}
