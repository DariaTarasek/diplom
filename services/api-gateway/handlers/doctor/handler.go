package doctor

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	DoctorClient *clients.DoctorClient
	AuthClient   *clients.AuthClient
}

func NewHandler(doctorClient *clients.DoctorClient, authClient *clients.AuthClient) *DoctorHandler {
	return &DoctorHandler{
		DoctorClient: doctorClient,
		AuthClient:   authClient,
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
	rg.GET("/doctor/consultation/patient-tests/:id", h.getPatientDocs)
	rg.GET("/doctor/consultation/patient-tests/download/:id", h.DownloadDocument)
	rg.GET("/doctor/me", h.getDoctorProfile)
	//  сюда остальные
}
