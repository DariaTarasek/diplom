package clients

import (
	statpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/statistics"
	"google.golang.org/grpc"
)

type StatisticsClient struct {
	Conn   *grpc.ClientConn
	Client statpb.StatisticsServiceClient
}

func NewStatisticsClient(address string) (*StatisticsClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := statpb.NewStatisticsServiceClient(conn)
	return &StatisticsClient{
		Conn:   conn,
		Client: client,
	}, nil
}
