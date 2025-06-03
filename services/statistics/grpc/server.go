package grpc

import (
	"context"
	pb "github.com/DariaTarasek/diplom/services/statistics/proto/statistics"
	"github.com/DariaTarasek/diplom/services/statistics/service"
)

type Server struct {
	pb.UnimplementedStatisticsServiceServer
	Service *service.StatisticsService
}

func (s *Server) GetAllStats(ctx context.Context, request *pb.EmptyRequest) (*pb.AllStatsResponse, error) {
	allStats, err := s.Service.GetAllStats(ctx)
	if err != nil {
		return &pb.AllStatsResponse{}, err
	}
	var topServices []*pb.TopService
	for _, item := range allStats.TopServices {
		topService := &pb.TopService{
			Name:  item.Name,
			Count: int32(item.Count),
		}
		topServices = append(topServices, topService)
	}
	var docAvgCheck []*pb.DoctorCheck
	for _, item := range allStats.DoctorCheckStat {
		docCheck := &pb.DoctorCheck{
			Doctor:   item.Doctor,
			AvgCheck: item.AvgCheck,
		}
		docAvgCheck = append(docAvgCheck, docCheck)
	}

	var docAvgVisit []*pb.DoctorAvgVisit
	for _, item := range allStats.DoctorAvgVisit {
		docVisit := &pb.DoctorAvgVisit{
			Doctor:          item.Doctor,
			AvgWeeklyVisits: item.AvgWeeklyVisits,
		}
		docAvgVisit = append(docAvgVisit, docVisit)
	}

	var docsPatients []*pb.DoctorUniquePatient
	for _, item := range allStats.DoctorUniquePatients {
		docPatient := &pb.DoctorUniquePatient{
			Doctor:         item.DoctorID,
			UniquePatients: int32(item.UniquePatients),
		}
		docsPatients = append(docsPatients, docPatient)
	}

	var ageGroups []*pb.AgeGroupStat
	for _, item := range allStats.AgeGroupStat {
		group := &pb.AgeGroupStat{
			AgeGroup: item.AgeGroup,
			Percent:  item.Percent,
		}
		ageGroups = append(ageGroups, group)
	}

	return &pb.AllStatsResponse{
		TotalPatients:        int32(allStats.TotalPatients),
		TotalVisits:          int32(allStats.TotalVisits),
		TopServices:          topServices,
		DocAvgVisits:         docAvgVisit,
		DoctorCheck:          docAvgCheck,
		DoctorPatients:       docsPatients,
		AgeGroups:            ageGroups,
		NewPatientsThisMonth: int32(allStats.NewPatientsThisMonth),
		AvgVisitPerPatient:   allStats.AvgVisitPerPatient,
		TotalIncome:          allStats.TotalIncome,
		MonthlyIncome:        allStats.MonthlyIncome,
		ClinicAvgCheck:       allStats.ClinicAvgCheck,
	}, nil
}
