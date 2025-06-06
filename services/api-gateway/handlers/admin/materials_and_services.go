package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// AddMaterial godoc
// @Summary Добавить новый материал
// @Tags Администратор
// @Description Добавляет новый материал в систему
// @Accept json
// @Produce json
// @Param material body model.Material true "Информация о материале"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/materials [post]
func (h *Handler) AddMaterial(c *gin.Context) {
	var materialReq model.Material
	if err := c.ShouldBindJSON(&materialReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	AddMaterialRequest := &adminpb.AddMaterialRequest{
		Name:  materialReq.Name,
		Price: int32(materialReq.Price),
	}

	_, err := h.AdminClient.Client.AddMaterial(c.Request.Context(), AddMaterialRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// AddService godoc
// @Summary Добавить новую услугу
// @Tags Администратор
// @Description Добавляет новую услугу в систему
// @Accept json
// @Produce json
// @Param service body model.Service true "Информация об услуге"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/services [post]
func (h *Handler) AddService(c *gin.Context) {
	var serviceReq model.Service
	if err := c.ShouldBindJSON(&serviceReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	AddServiceRequest := &adminpb.AddServiceRequest{
		Name:  serviceReq.Name,
		Price: int32(serviceReq.Price),
		Type:  int32(serviceReq.ServiceTypeId),
	}

	_, err := h.AdminClient.Client.AddService(c.Request.Context(), AddServiceRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateMaterial godoc
// @Summary Обновить материал
// @Tags Администратор
// @Description Обновляет данные материала по ID
// @Accept json
// @Produce json
// @Param id path int true "ID материала"
// @Param material body model.Material true "Новые данные материала"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/materials/{id} [put]
func (h *Handler) UpdateMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var materialReq model.Material
	if err := c.ShouldBindJSON(&materialReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdateMaterialRequest := &adminpb.UpdateMaterialRequest{
		Id:    int32(id),
		Name:  materialReq.Name,
		Price: int32(materialReq.Price),
	}

	_, err = h.AdminClient.Client.UpdateMaterial(c.Request.Context(), UpdateMaterialRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateService godoc
// @Summary Обновить услугу
// @Tags Администратор
// @Description Обновляет данные услуги по ID
// @Accept json
// @Produce json
// @Param id path int true "ID услуги"
// @Param service body model.Service true "Новые данные услуги"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	var serviceReq model.Service
	if err := c.ShouldBindJSON(&serviceReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdateServiceRequest := &adminpb.UpdateServiceRequest{
		Id:    int32(id),
		Name:  serviceReq.Name,
		Price: int32(serviceReq.Price),
		Type:  int32(serviceReq.ServiceTypeId),
	}

	_, err = h.AdminClient.Client.UpdateService(c.Request.Context(), UpdateServiceRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteMaterial godoc
// @Summary Удалить материал
// @Tags Администратор
// @Description Удаляет материал по ID
// @Produce json
// @Param id path int true "ID материала"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/materials/{id} [delete]
func (h *Handler) DeleteMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	DeleteRequest := &adminpb.DeleteRequest{Id: int32(id)}

	_, err = h.AdminClient.Client.DeleteMaterial(c.Request.Context(), DeleteRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteService godoc
// @Summary Удалить услугу
// @Tags Администратор
// @Description Удаляет услугу по ID
// @Produce json
// @Param id path int true "ID услуги"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/services/{id} [delete]

func (h *Handler) DeleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	DeleteRequest := &adminpb.DeleteRequest{Id: int32(id)}

	_, err = h.AdminClient.Client.DeleteService(c.Request.Context(), DeleteRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
