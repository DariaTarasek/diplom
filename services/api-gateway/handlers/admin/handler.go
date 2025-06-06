package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	AdminClient      *clients.AdminClient
	AuthClient       *clients.AuthClient
	AccessMiddleware func(requiredPermission int32) gin.HandlerFunc
}

func NewHandler(adminClient *clients.AdminClient, authClient *clients.AuthClient, accessMiddleware func(requiredPermission int32) gin.HandlerFunc) *Handler {
	return &Handler{
		AdminClient:      adminClient,
		AuthClient:       authClient,
		AccessMiddleware: accessMiddleware,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/clinic-schedule", h.AccessMiddleware(5), h.UpdateClinicSchedule)
	//rg.GET("/admin-data", h.GetUserRole)
	rg.POST("/doctor-schedule/:selectedDoctor", h.AccessMiddleware(5), h.UpdateDoctorSchedule)
	rg.POST("/clinic-overrides", h.AccessMiddleware(7), h.AddClinicDailyOverride)
	rg.POST("/doctor-overrides", h.AccessMiddleware(7), h.AddDoctorDailyOverride)
	rg.POST("/materials", h.AccessMiddleware(18), h.AddMaterial)
	rg.POST("/services", h.AccessMiddleware(18), h.AddService)
	rg.PUT("/materials/:id", h.AccessMiddleware(18), h.UpdateMaterial)
	rg.PUT("/services/:id", h.AccessMiddleware(18), h.UpdateService)
	rg.DELETE("/materials/:id", h.AccessMiddleware(18), h.DeleteMaterial)
	rg.DELETE("/services/:id", h.AccessMiddleware(18), h.DeleteService)
	rg.GET("/staff-admins", h.AccessMiddleware(1), h.GetAdmins)
	rg.GET("/staff-doctors", h.AccessMiddleware(1), h.GetDoctors)
	rg.GET("/patients", h.AccessMiddleware(1), h.GetPatients)
	rg.POST("/save-admin", h.AccessMiddleware(9), h.UpdateAdmin)
	rg.POST("save-doctor", h.AccessMiddleware(9), h.UpdateDoctor)
	rg.PUT("/patients/:id", h.AccessMiddleware(8), h.UpdatePatient)
	rg.DELETE("/patients/:id", h.AccessMiddleware(19), h.DeleteUser)
	rg.DELETE("/admins/:id", h.AccessMiddleware(19), h.DeleteUser)
	rg.DELETE("/doctors/:id", h.AccessMiddleware(19), h.DeleteUser)
	rg.PUT("/admins-login/:id", h.AccessMiddleware(9), h.UpdateEmployeeLogin)
	rg.PUT("/doctors-login/:id", h.AccessMiddleware(9), h.UpdateEmployeeLogin)
	rg.PUT("/patients-login/:id", h.AccessMiddleware(8), h.UpdatePatientLogin)
	rg.GET("/completed-visits", h.AccessMiddleware(1), h.GetVisitPayments)
	rg.GET("/schedule-admin", h.AccessMiddleware(1), h.GetScheduleGrid)
	rg.GET("/admin/me", h.AccessMiddleware(1), h.getAdminProfile)
	rg.GET("/unconfirmed-appointments", h.AccessMiddleware(1), h.GetUnconfirmedAppointments)
	rg.PUT("/completed-visits/:id", h.AccessMiddleware(10), h.UpdateVisitPayment)
	rg.PUT("/unconfirmed-appointments/:id", h.AccessMiddleware(14), h.UpdateAppointment)

	//	rg.GET("/specialties", h.GetSpecs)
}
