package clients

import (
	patientpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/patient"
	"google.golang.org/grpc"
)

type PatientClient struct {
	Conn   *grpc.ClientConn
	Client patientpb.PatientServiceClient
}

func NewPatientClient(address string) (*PatientClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := patientpb.NewPatientServiceClient(conn)
	return &PatientClient{
		Conn:   conn,
		Client: client,
	}, nil
}
