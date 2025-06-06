package auth

import (
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type recoveryRequest struct {
	Login string `json:"login"`
}

// @Summary Восстановление пароля сотрудника
// @Accept json
// @Produce json
// @Tags Авторизация
// @Param input body recoveryRequest true "Логин"
// @Success 200 {object} gin.H
// @Failure 400,500 {object} gin.H
// @Router /api/employee-password-recovery [post]
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

// @Summary Восстановление пароля пациента
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param input body recoveryRequest true "Логин"
// @Success 200 {object} gin.H
// @Failure 400,500 {object} gin.H
// @Router /api/password-recovery [post]
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
