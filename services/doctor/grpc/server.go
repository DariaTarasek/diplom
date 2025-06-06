package grpcserver

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/doctor/model"
	pb "github.com/DariaTarasek/diplom/services/doctor/proto/doctor"
	"github.com/DariaTarasek/diplom/services/doctor/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedDoctorServiceServer
	Service *service.DoctorService
}

func (s *Server) GetTodayAppointments(ctx context.Context, request *pb.GetTodayAppointmentsRequest) (*pb.GetTodayAppointmentsResponse, error) {
	apps, err := s.Service.GetTodayAppointments(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	todays := make([]*pb.TodayAppointment, 0, len(apps))
	for _, app := range apps {
		today := &pb.TodayAppointment{
			Id:        int32(app.ID),
			Date:      app.Date,
			Time:      app.Time,
			PatientID: int32(app.PatientID),
			Patient:   app.Patient,
		}
		todays = append(todays, today)
	}
	return &pb.GetTodayAppointmentsResponse{Appointments: todays}, nil
}

func (s *Server) GetUpcomingAppointments(ctx context.Context, request *pb.GetUpcomingAppointmentsRequest) (*pb.GetUpcomingAppointmentsResponse, error) {
	apps, err := s.Service.GetUpcomingAppointments(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	var cells []*pb.ScheduleCell
	for date, row := range apps.Table {
		for timeStr, appointment := range row {
			cell := &pb.ScheduleCell{
				Date: date,
				Time: timeStr,
			}
			if appointment != nil {
				cell.Appointment = &pb.UpcomingAppointment{
					Id:        int64(appointment.ID),
					PatientId: int64(appointment.PatientID),
					Patient:   appointment.Patient,
				}
			}
			cells = append(cells, cell)
		}
	}

	resp := &pb.ScheduleTable{
		Dates: apps.Dates,
		Times: apps.Times,
		Table: cells,
	}

	return &pb.GetUpcomingAppointmentsResponse{Schedule: resp}, nil
}

func (s *Server) GetPatientAllergiesChronics(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetPatientAllergiesChronicsResponse, error) {
	notes, err := s.Service.GetPatientAllergiesChronics(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	var gRPCNotes []*pb.PatientAllergiesChronics
	for _, item := range notes {
		note := &pb.PatientAllergiesChronics{
			Id:        int32(item.ID),
			PatientId: int32(item.PatientID),
			Type:      item.Type,
			Title:     item.Title,
		}
		gRPCNotes = append(gRPCNotes, note)
	}
	return &pb.GetPatientAllergiesChronicsResponse{
		PatientAllergiesChronics: gRPCNotes}, nil
}

func (s *Server) GetAppointmentByID(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetAppointmentByIDResponse, error) {
	appt, err := s.Service.GetAppointmentByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	appointment := &pb.Appointment{
		Id:          int32(appt.ID),
		DoctorId:    int32(appt.DoctorID),
		Date:        timestamppb.New(appt.Date),
		Time:        timestamppb.New(appt.Time),
		PatientId:   int32(*appt.PatientID),
		SecondName:  appt.PatientSecondName,
		FirstName:   appt.PatientFirstName,
		Surname:     *appt.PatientSurname,
		BirthDate:   timestamppb.New(appt.PatientBirthDate),
		Gender:      appt.PatientGender,
		PhoneNumber: appt.PatientPhoneNumber,
		Status:      appt.Status,
		CreatedAt:   timestamppb.New(appt.CreatedAt),
		UpdatedAt:   timestamppb.New(appt.UpdatedAt),
	}
	return &pb.GetAppointmentByIDResponse{Appt: appointment}, nil
}

func ConvertToPointerSlice(input []model.Diagnose) []*pb.Diagnose {
	result := make([]*pb.Diagnose, 0, len(input))
	for _, d := range input {
		diag := pb.Diagnose{
			IcdCode: d.ICDCode,
			Notes:   d.Notes,
		}
		result = append(result, &diag)
	}
	return result
}

func (s *Server) GetPatientVisits(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetPatientVisitsResponse, error) {
	visits, err := s.Service.GetPatientVisits(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	var gRPCVisits []*pb.Visit
	for _, item := range visits {
		grpcVisit := &pb.Visit{
			Id:         int32(item.ID),
			ApptId:     int32(item.AppointmentID),
			PatientId:  int32(item.PatientID),
			Doctor:     item.Doctor,
			Complaints: item.Complaints,
			Treatment:  item.Treatment,
			CreatedAt:  item.CreatedAt,
			Diagnoses:  ConvertToPointerSlice(item.Diagnoses),
		}
		gRPCVisits = append(gRPCVisits, grpcVisit)
	}

	return &pb.GetPatientVisitsResponse{Visits: gRPCVisits}, nil
}

func (s *Server) AddPatientAllergiesChronics(ctx context.Context, req *pb.AddPatientAllergiesChronicsRequest) (*pb.DefaultResponse, error) {
	var notes []model.AllergiesChronics
	for _, item := range req.Notes {
		note := model.AllergiesChronics{
			PatientID: int(item.PatientId),
			Type:      item.Type,
			Title:     item.Title,
		}
		notes = append(notes, note)
	}
	err := s.Service.AddPatientAllergiesChronics(ctx, notes)
	if err != nil {
		return nil, fmt.Errorf("не удалось добавить мед. запись: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddConsultation(ctx context.Context, req *pb.AddConsultationRequest) (*pb.AddConsultationResponse, error) {
	visit := model.AddVisit{
		AppointmentID: int(req.AppointmentId),
		PatientID:     int(req.PatientId),
		DoctorID:      int(req.DoctorId),
		Complaints:    req.Complaints,
		Treatment:     req.Treatment,
	}

	var materials []model.VisitMaterial
	for _, m := range req.Materials {
		materials = append(materials, model.VisitMaterial{
			VisitID:    int(m.VisitId),
			MaterialID: int(m.MaterialId),
			Amount:     int(m.Amount),
		})
	}

	var services []model.VisitService
	for _, s := range req.Services {
		services = append(services, model.VisitService{
			VisitID:   int(s.VisitId),
			ServiceID: int(s.ServiceId),
			Amount:    int(s.Amount),
		})
	}

	var diagnoses []model.VisitDiagnose
	for _, d := range req.Diagnoses {
		diagnoses = append(diagnoses, model.VisitDiagnose{
			VisitID:   int(d.VisitId),
			ICDCodeID: int(d.IcdCodeId),
			Note:      d.Note,
		})
	}

	err := s.Service.AddConsultation(ctx, materials, services, diagnoses, visit, req.Token)
	if err != nil {
		return nil, fmt.Errorf("не удалось провести приём: %w", err)
	}
	return &pb.AddConsultationResponse{}, nil
}

func (s *Server) AddVisitPayment(ctx context.Context, req *pb.VisitPaymentRequest) (*pb.DefaultResponse, error) {
	err := s.Service.AddVisitPayment(ctx, model.VisitPayment{
		VisitID: int(req.VisitId),
		Price:   int(req.Price),
		Status:  req.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось добавить платеж по визиту: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateVisitPayment(ctx context.Context, req *pb.VisitPaymentRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdateVisitPayment(ctx, model.VisitPayment{
		VisitID: int(req.VisitId),
		Price:   int(req.Price),
		Status:  req.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить платеж по визиту: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetDocumentsByPatientID(ctx context.Context, request *pb.GetDocumentsRequest) (*pb.GetDocumentsResponse, error) {
	resp, err := s.Service.GetDocumentsInfo(ctx, int(request.PatientID))
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
