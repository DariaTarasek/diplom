package auth

import (
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type recoveryRequest struct {
	Login string `json:"login"`
}

func (h *Handler) EmployeePasswordRecovery(c *gin.Context) {
	var req recoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.AuthClient.Client.EmployeePasswordRecovery(c.Request.Context(), &authpb.EmployeePasswordRecoveryRequest{Login: req.Login})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пароль изменен!"})
}

func (h *Handler) PatientPasswordRecovery(c *gin.Context) {
	var req recoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.AuthClient.Client.PatientPasswordRecovery(c.Request.Context(), &authpb.PatientPasswordRecoveryRequest{Login: req.Login})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пароль изменен!"})
}
