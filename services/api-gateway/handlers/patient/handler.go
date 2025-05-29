package patient

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	PatientClient *clients.PatientClient
}

func NewHandler(patientClient *clients.PatientClient) *PatientHandler {
	return &PatientHandler{
		PatientClient: patientClient,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *PatientHandler) {
	rg.GET("/appointment-doctor-schedule/:doctorId", h.getAppointmentSlots)
	rg.POST("/appointments", h.addAppointment)
	//  сюда остальные
}
