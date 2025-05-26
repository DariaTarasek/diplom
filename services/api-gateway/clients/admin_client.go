package clients

import (
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"google.golang.org/grpc"
)

type AdminClient struct {
	Conn   *grpc.ClientConn
	Client adminpb.AdminServiceClient
}

func NewAdminClient(address string) (*AdminClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := adminpb.NewAdminServiceClient(conn)
	return &AdminClient{
		Conn:   conn,
		Client: client,
	}, nil
}
