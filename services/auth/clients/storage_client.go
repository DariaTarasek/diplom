package clients

import (
	"github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"google.golang.org/grpc"
)

type StorageClient struct {
	Conn   *grpc.ClientConn
	Client storagepb.StorageServiceClient
}

func NewStorageClient(address string) (*StorageClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure()) // или grpc.WithTransportCredentials
	if err != nil {
		return nil, err
	}

	client := storagepb.NewStorageServiceClient(conn)
	return &StorageClient{
		Conn:   conn,
		Client: client,
	}, nil
}
