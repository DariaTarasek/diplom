package doctor

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	DoctorClient     *clients.DoctorClient
	AuthClient       *clients.AuthClient
	AccessMiddleware func(requiredPermission int32) gin.HandlerFunc
}

func NewHandler(doctorClient *clients.DoctorClient, authClient *clients.AuthClient, accessMiddleware func(requiredPermission int32) gin.HandlerFunc) *DoctorHandler {
	return &DoctorHandler{
		DoctorClient:     doctorClient,
		AuthClient:       authClient,
		AccessMiddleware: accessMiddleware,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *DoctorHandler) {
	rg.GET("/appointments-today", h.AccessMiddleware(2), h.GetTodayAppointments)
	rg.GET("/schedule-with-appointments", h.AccessMiddleware(2), h.GetUpcomingAppointments)
	rg.GET("/patient-notes/:id", h.AccessMiddleware(11), h.GetPatientAllergiesChronics)
	rg.GET("/appointments/:id", h.AccessMiddleware(2), h.GetAppointmentByID)
	rg.GET("/patient-history/:id", h.AccessMiddleware(12), h.GetPatientVisits)
	rg.POST("/visits", h.AccessMiddleware(13), h.AddConsultation)
	rg.POST("/patient-notes/:id", h.AccessMiddleware(11), h.AddPatientAllergiesChronics)
	rg.GET("/doctor/consultation/patient-tests/:id", h.AccessMiddleware(2), h.getPatientDocs)
	rg.GET("/doctor/consultation/patient-tests/download/:id", h.AccessMiddleware(16), h.DownloadDocument)
	rg.GET("/doctor/me", h.AccessMiddleware(2), h.getDoctorProfile)
	//  сюда остальные
}
