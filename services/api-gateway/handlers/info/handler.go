package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	store *clients.StorageClient
}

func NewInfoHandler(store *clients.StorageClient) *InfoHandler {
	return &InfoHandler{store: store}
}

func RegisterRoutes(rg *gin.RouterGroup, h *InfoHandler) {
	rg.GET("/specialties", h.GetAllSpecs)
	rg.GET("/doctors", h.GetDoctors)
	rg.GET("/clinic-schedule", h.GetClinicWeeklySchedule)
	rg.GET("/doctors/:specialty", h.getDoctorsBySpecialty)
	// добавляй сюда остальные
}
