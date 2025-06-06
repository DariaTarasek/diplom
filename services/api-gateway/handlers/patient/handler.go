package patient

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	PatientClient    *clients.PatientClient
	AccessMiddleware func(requiredPermission int32) gin.HandlerFunc
}

func NewHandler(patientClient *clients.PatientClient, accessMiddleware func(requiredPermission int32) gin.HandlerFunc) *PatientHandler {
	return &PatientHandler{
		PatientClient:    patientClient,
		AccessMiddleware: accessMiddleware,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *PatientHandler) {
	rg.GET("/appointment-doctor-schedule/:doctorId", h.getAppointmentSlots)
	rg.POST("/appointments", h.addAppointment)
	rg.GET("/patient/upcoming", h.AccessMiddleware(3), h.getUpcomingAppointments)
	rg.GET("/patient/history", h.AccessMiddleware(3), h.getHistoryVisits)
	rg.GET("/patient/tests", h.AccessMiddleware(16), h.getDocuments)
	rg.PUT("/appointments/transfer", h.AccessMiddleware(14), h.UpdateAppointment)
	rg.GET("/appointments/cancel/:id", h.AccessMiddleware(14), h.CancelAppointment)
	rg.POST("/patient/tests/upload", h.AccessMiddleware(15), h.UploadTest)
	rg.GET("/patient/tests/:id/download", h.AccessMiddleware(16), h.DownloadDocument)
}
