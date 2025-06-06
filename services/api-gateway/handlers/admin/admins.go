package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetAdmins godoc
// @Summary Получить список администраторов
// @Tags Администратор
// @Description Возвращает список всех администраторов
// @Produce json
// @Success 200 {array} model.AdminForAdminList
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/staff-admin [get]
func (h *Handler) GetAdmins(c *gin.Context) {
	items, err := h.AdminClient.Client.GetAdmins(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var admins []model.AdminForAdminList
	for _, item := range items.Admins {
		admin := model.AdminForAdminList{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Gender:      item.Gender,
			Role:        item.Role,
		}
		admins = append(admins, admin)
	}
	c.JSON(http.StatusOK, admins)
}

// UpdateAdmin godoc
// @Summary Обновить данные администратора
// @Tags Администратор
// @Description Обновляет информацию об администраторе
// @Accept json
// @Produce json
// @Param admin body model.AdminForAdminList true "Данные администратора"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/save-admin [put]
func (h *Handler) UpdateAdmin(c *gin.Context) {
	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	log.Println(err.Error())
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
	//	return
	//}

	var adminReq model.AdminForAdminList
	if err := c.ShouldBindJSON(&adminReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdateAdminRequest := &adminpb.UpdateAdminRequest{
		UserId:      int32(adminReq.ID),
		FirstName:   adminReq.FirstName,
		SecondName:  adminReq.SecondName,
		Surname:     adminReq.Surname,
		PhoneNumber: adminReq.PhoneNumber,
		Email:       adminReq.Email,
		Gender:      adminReq.Gender,
		Role:        adminReq.Role,
	}

	_, err := h.AdminClient.Client.UpdateAdmin(c.Request.Context(), UpdateAdminRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
