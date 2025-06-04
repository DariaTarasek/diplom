package clients

import (
	storagepb "github.com/DariaTarasek/diplom/services/doctor/proto/storage"
	"google.golang.org/grpc"
)

type StorageClient struct {
	Conn   *grpc.ClientConn
	Client storagepb.StorageServiceClient
}

func NewStorageClient(address string) (*StorageClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50*1024*1024), grpc.MaxCallSendMsgSize(50*1024*1024))) // или grpc.WithTransportCredentials
	if err != nil {
		return nil, err
	}

	client := storagepb.NewStorageServiceClient(conn)
	return &StorageClient{
		Conn:   conn,
		Client: client,
	}, nil
}
