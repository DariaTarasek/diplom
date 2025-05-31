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
	//  сюда остальные
}
