package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	AdminClient *clients.AdminClient
}

func NewHandler(adminClient *clients.AdminClient) *Handler {
	return &Handler{
		AdminClient: adminClient,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/clinic-schedule", h.UpdateClinicSchedule)
	rg.GET("/admin-data", h.GetUserRole)
	rg.POST("/doctor-schedule/:selectedDoctor", h.UpdateDoctorSchedule)
	rg.POST("/clinic-overrides", h.AddClinicDailyOverride)
	rg.POST("/doctor-overrides", h.AddDoctorDailyOverride)
	rg.POST("/materials", h.AddMaterial)
	rg.POST("/services", h.AddService)
	rg.PUT("/materials/:id", h.UpdateMaterial)
	rg.PUT("/services/:id", h.UpdateService)
	rg.DELETE("/materials/:id", h.DeleteMaterial)
	rg.DELETE("/services/:id", h.DeleteService)
	rg.GET("/staff-admins", h.GetAdmins)
	rg.GET("/staff-doctors", h.GetDoctors)
	rg.GET("/patients", h.GetPatients)
	rg.POST("/save-admin", h.UpdateAdmin)
	rg.POST("save-doctor", h.UpdateDoctor)
	rg.PUT("/patients/:id", h.UpdatePatient)
	rg.DELETE("/patients/:id", h.DeleteUser)
	rg.DELETE("/admins/:id", h.DeleteUser)
	rg.DELETE("/doctors/:id", h.DeleteUser)
	rg.PUT("/admins-login/:id", h.UpdateEmployeeLogin)
	rg.PUT("/doctors-login/:id", h.UpdateEmployeeLogin)
	rg.PUT("/patients-login/:id", h.UpdatePatientLogin)
	rg.GET("/completed-visits", h.GetVisitPayments)
	//	rg.GET("/specialties", h.GetSpecs)
}
