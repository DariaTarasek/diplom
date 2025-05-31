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
	rg.POST("/employee-register", h.EmployeeRegister)
	rg.POST("/register", h.PatientRegister)
	rg.POST("/register-in-clinic", h.PatientRegisterInClinic)
	rg.POST("/employee-password-recovery", h.EmployeePasswordRecovery)
	rg.POST("/password-recovery", h.PatientPasswordRecovery)
	rg.POST("/request-code", h.requestCode)
	rg.POST("/verify-code", h.verifyCode)
	rg.POST("/login", h.authorize)
	rg.GET("/patient/me", h.getPatient)
	//  сюда остальные
}
