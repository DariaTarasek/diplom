package auth

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthClient       *clients.AuthClient
	AccessMiddleware func(requiredPermission int32) gin.HandlerFunc
}

func NewHandler(authClient *clients.AuthClient, accessMiddleware func(requiredPermission int32) gin.HandlerFunc) *Handler {
	return &Handler{
		AuthClient:       authClient,
		AccessMiddleware: accessMiddleware,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/employee-register", h.AccessMiddleware(9), h.EmployeeRegister)
	rg.POST("/register", h.PatientRegister)
	rg.POST("/register-in-clinic", h.AccessMiddleware(8), h.PatientRegisterInClinic)
	rg.POST("/employee-password-recovery", h.EmployeePasswordRecovery)
	rg.POST("/password-recovery", h.PatientPasswordRecovery)
	rg.POST("/request-code", h.requestCode)
	rg.POST("/verify-code", h.verifyCode)
	rg.POST("/login", h.authorize)
	rg.GET("/patient/me", h.AccessMiddleware(3), h.getPatient)
	//  сюда остальные
}
