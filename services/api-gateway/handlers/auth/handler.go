package auth

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthClient *clients.AuthClient
}

func NewHandler(authClient *clients.AuthClient) *Handler {
	return &Handler{
		AuthClient: authClient,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/employee-auth", h.EmployeeRegister)
	rg.POST("/auth", h.PatientRegister)
	rg.POST("/auth-in-clinic", h.PatientRegisterInClinic)
	rg.POST("/employee-password-recovery", h.EmployeePasswordRecovery)
	rg.POST("/password-recovery", h.PatientPasswordRecovery)
	//  сюда остальные
}
