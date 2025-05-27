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
}
