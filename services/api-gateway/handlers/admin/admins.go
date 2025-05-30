package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
