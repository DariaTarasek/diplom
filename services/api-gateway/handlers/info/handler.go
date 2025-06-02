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
	rg.GET("/doctor-schedule/:selectedDoctor", h.GetDoctorWeeklySchedule)
	rg.GET("doctor-overrides/:doctor_id/:date", h.GetDoctorOverride)
	rg.GET("clinic-overrides/:date", h.GetClinicOverride)
	rg.GET("/doctors/:specialty", h.GetDoctorsBySpecialty)
	rg.GET("/services", h.GetServices)
	rg.GET("/materials", h.GetMaterials)
	rg.GET("/icd-codes", h.GetICDCodes)
	rg.GET("/service-categories", h.GetServicesTypes)
	//rg.GET("/patient-history/:id", h.GetPatientHistory)
	//rg.GET("/patient-notes/:id", h.GetPatientNotes)
}
