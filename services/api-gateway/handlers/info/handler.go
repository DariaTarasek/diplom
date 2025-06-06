package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	store            *clients.StorageClient
	AccessMiddleware func(requiredPermission int32) gin.HandlerFunc
}

func NewInfoHandler(store *clients.StorageClient, accessMiddleware func(requiredPermission int32) gin.HandlerFunc) *InfoHandler {
	return &InfoHandler{
		store:            store,
		AccessMiddleware: accessMiddleware,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *InfoHandler) {
	rg.GET("/specialties", h.GetAllSpecs)
	rg.GET("/doctors", h.GetDoctors)
	rg.GET("/clinic-schedule", h.AccessMiddleware(1), h.GetClinicWeeklySchedule)
	rg.GET("/doctor-schedule/:selectedDoctor", h.AccessMiddleware(1), h.GetDoctorWeeklySchedule)
	rg.GET("/doctor-overrides/:doctor_id/:date", h.AccessMiddleware(1), h.GetDoctorOverride)
	rg.GET("/clinic-overrides/:date", h.AccessMiddleware(1), h.GetClinicOverride)
	rg.GET("/doctors/:specialty", h.GetDoctorsBySpecialty)
	rg.GET("/services", h.AccessMiddleware(20), h.GetServices)
	rg.GET("/materials", h.AccessMiddleware(20), h.GetMaterials)
	rg.GET("/icd-codes", h.AccessMiddleware(2), h.GetICDCodes)
	rg.GET("/service-categories", h.GetServicesTypes)
	//rg.GET("/patient-history/:id", h.GetPatientHistory)
	//rg.GET("/patient-notes/:id", h.GetPatientNotes)
}
