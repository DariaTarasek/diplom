package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

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
