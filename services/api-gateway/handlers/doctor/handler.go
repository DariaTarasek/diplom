package doctor

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	DoctorClient *clients.DoctorClient
}

func NewHandler(doctorClient *clients.DoctorClient) *DoctorHandler {
	return &DoctorHandler{
		DoctorClient: doctorClient,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *DoctorHandler) {
	rg.GET("/appointments-today", h.GetTodayAppointments)
	rg.GET("/schedule-with-appointments", h.GetUpcomingAppointments)
	rg.GET("/patient-notes/:id", h.GetPatientAllergiesChronics)
	rg.GET("/appointments/:id", h.GetAppointmentByID)
	rg.GET("/patient-history/:id", h.GetPatientVisits)
	rg.POST("/visits", h.AddConsultation)
	rg.POST("/patient-notes/:id", h.AddPatientAllergiesChronics)
	//  сюда остальные
}
