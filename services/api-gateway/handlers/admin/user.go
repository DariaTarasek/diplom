package admin

import (
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type employeeLoginRequest struct {
	Email string `json:"email"`
}

type patientLoginRequest struct {
	Phone string `json:"phone"`
}

// @Summary Удалить пользователя
// @Tags Администратор
// @Description Удаляет пользователя по ID
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/doctors/{id} [delete]
// @Router /api/admins/{id} [delete]
// @Router /api/patients/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	DeleteRequest := &adminpb.DeleteRequest{Id: int32(id)}

	_, err = h.AdminClient.Client.DeleteUser(c.Request.Context(), DeleteRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) UpdateEmployeeLogin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var loginReq employeeLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdateEmployeeLoginRequest := &adminpb.UpdateUserLoginRequest{
		UserId: int32(id),
		Login:  loginReq.Email,
	}
	_, err = h.AdminClient.Client.UpdateEmployeeLogin(c.Request.Context(), UpdateEmployeeLoginRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Обновить логин сотрудника
// @Tags Администратор
// @Description Обновляет логин сотрудника по ID
// @Accept json
// @Produce json
// @Param id path int true "ID сотрудника"
// @Param login body employeeLoginRequest true "Новый логин"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/patients-login/{id} [put]
func (h *Handler) UpdatePatientLogin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var loginReq patientLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdatePatientLoginRequest := &adminpb.UpdateUserLoginRequest{
		UserId: int32(id),
		Login:  loginReq.Phone,
	}
	_, err = h.AdminClient.Client.UpdatePatientLogin(c.Request.Context(), UpdatePatientLoginRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
