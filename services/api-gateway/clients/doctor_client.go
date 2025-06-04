package clients

import (
	doctorpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/doctor"
	"google.golang.org/grpc"
)

type DoctorClient struct {
	Conn   *grpc.ClientConn
	Client doctorpb.DoctorServiceClient
}

func NewDoctorClient(address string) (*DoctorClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50*1024*1024), grpc.MaxCallSendMsgSize(50*1024*1024))) // или grpc.WithTransportCredentials
	if err != nil {
		return nil, err
	}

	client := doctorpb.NewDoctorServiceClient(conn)
	return &DoctorClient{
		Conn:   conn,
		Client: client,
	}, nil
}
