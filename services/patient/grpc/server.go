package grpcserver

import (
	"context"
	"github.com/DariaTarasek/diplom/services/patient/model"
	pb "github.com/DariaTarasek/diplom/services/patient/proto/patient"
	"github.com/DariaTarasek/diplom/services/patient/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedPatientServiceServer
	Service *service.PatientService
}

func (s *Server) GetAppointmentSlots(ctx context.Context, request *pb.GetAppointmentSlotsRequest) (*pb.GetAppointmentSlotsResponse, error) {
	slots, err := s.Service.MakeDoctorAppointmentSlots(ctx, int(request.DoctorId))
	if err != nil {
		return nil, err
	}
	var res []*pb.DaySlots
	for _, slot := range slots {
		res = append(res, &pb.DaySlots{
			Date:  slot.Label,
			Slots: slot.Slots,
		})
	}
	return &pb.GetAppointmentSlotsResponse{Slots: res}, nil
}

func (s *Server) AddAppointment(ctx context.Context, request *pb.AddAppointmentRequest) (*pb.DefaultResponse, error) {
	patientID := model.UserID(request.Appointment.PatientId)
	appointment := model.Appointment{
		DoctorID:           model.UserID(request.Appointment.DoctorId),
		PatientID:          &patientID,
		Date:               request.Appointment.Date.AsTime(),
		Time:               request.Appointment.Time.AsTime(),
		PatientSecondName:  request.Appointment.SecondName,
		PatientFirstName:   request.Appointment.FirstName,
		PatientSurname:     &request.Appointment.Surname,
		PatientBirthDate:   request.Appointment.BirthDate.AsTime(),
		PatientGender:      request.Appointment.Gender,
		PatientPhoneNumber: request.Appointment.PhoneNumber,
	}
	err := s.Service.AddAppointment(ctx, appointment)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetUpcomingAppointments(ctx context.Context, request *pb.GetUpcomingAppointmentsRequest) (*pb.GetUpcomingAppointmentsResponse, error) {
	apps, err := s.Service.GetUpcomingAppointments(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	upcoming := make([]*pb.UpcomingAppointments, 0, len(apps))
	for _, app := range apps {
		appointment := &pb.UpcomingAppointments{
			Id:        int32(app.ID),
			Date:      app.Date,
			Time:      app.Time,
			DoctorId:  int32(app.DoctorID),
			Doctor:    app.Doctor,
			Specialty: app.Specialty,
		}
		upcoming = append(upcoming, appointment)
	}
	return &pb.GetUpcomingAppointmentsResponse{Appointments: upcoming}, nil
}

func (s *Server) UpdateAppointment(ctx context.Context, request *pb.UpdateAppointmentRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdateAppointment(ctx, model.Appointment{
		ID:   model.AppointmentID(request.Appointment.Id),
		Date: request.Appointment.Date.AsTime(),
		Time: request.Appointment.Time.AsTime(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) CancelAppointment(ctx context.Context, request *pb.GetByIDRequest) (*pb.DefaultResponse, error) {
	err := s.Service.CancelAppointment(ctx, model.AppointmentID(request.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetHistoryVisits(ctx context.Context, request *pb.GetHistoryVisitsRequest) (*pb.GetHistoryVisitsResponse, error) {
	resp, err := s.Service.GetHistoryVisits(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	historyVisits := make([]*pb.HistoryVisit, 0, len(resp))
	for _, item := range resp {
		historyVisit := &pb.HistoryVisit{
			Id:        int32(item.ID),
			Date:      item.Date,
			DoctorId:  int32(item.DoctorID),
			Doctor:    item.Doctor,
			Diagnose:  item.Diagnose,
			Treatment: item.Treatment,
		}
		historyVisits = append(historyVisits, historyVisit)
	}
	return &pb.GetHistoryVisitsResponse{Visits: historyVisits}, nil
}

func (s *Server) UploadTest(ctx context.Context, req *pb.UploadTestRequest) (*pb.UploadTestResponse, error) {
	docID, err := s.Service.UploadTest(ctx, model.UploadTestInput{
		Token:       req.Token,
		FileName:    req.FileName,
		FileContent: req.FileContent,
		Description: req.Description,
	})
	if err != nil {
		// маппинг ошибок в gRPC status
		return nil, status.Errorf(codes.InvalidArgument, "upload failed: %v", err)
	}

	return &pb.UploadTestResponse{
		DocumentId: docID.String(),
	}, nil
}

func (s *Server) GetDocumentsByPatientID(ctx context.Context, request *pb.GetDocumentsRequest) (*pb.GetDocumentsResponse, error) {
	resp, err := s.Service.GetDocumentsInfo(ctx, request.Token)
	if err != nil {
		return nil, err
	}

	docs := make([]*pb.DocumentInfo, 0, len(resp))
	for _, item := range resp {
		doc := &pb.DocumentInfo{
			Id:          item.ID,
			FileName:    item.FileName,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		}
		docs = append(docs, doc)
	}
	return &pb.GetDocumentsResponse{Documents: docs}, nil
}

func (s *Server) DownloadDocument(ctx context.Context, request *pb.DownloadDocumentRequest) (*pb.DownloadDocumentResponse, error) {
	resp, err := s.Service.DownloadDocument(ctx, request.DocumentId)
	if err != nil {
		return nil, err
	}
	return &pb.DownloadDocumentResponse{
		FileName:    resp.FileName,
		FileContent: resp.FileContent,
	}, nil
}
