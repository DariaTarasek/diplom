package info

import (
	"fmt"
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *InfoHandler) GetServices(c *gin.Context) {
	items, err := h.store.Client.GetServices(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var services []model.Service
	for _, item := range items.Services {
		service := model.Service{
			ID:            int(item.Id),
			Name:          item.Name,
			Price:         int(item.Price),
			ServiceTypeId: int(item.Type),
		}
		services = append(services, service)
	}
	c.JSON(http.StatusOK, gin.H{"services": services})
}

func (h *InfoHandler) GetServicesTypes(c *gin.Context) {
	items, err := h.store.Client.GetServicesTypes(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var types []model.ServiceType
	for _, item := range items.Types {
		sType := model.ServiceType{
			ID:   int(item.Id),
			Name: item.Name,
		}
		types = append(types, sType)
	}
	c.JSON(http.StatusOK, gin.H{"categories": types})
}

func (h *InfoHandler) GetServiceTypeById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный запрос " + err.Error()})
		return
	}

	fmt.Println(id)
	item, err := h.store.Client.GetServiceTypeById(c.Request.Context(), &storagepb.GetServiceTypeByIdRequest{Id: int32(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список врачей по специальности " + err.Error()})
		return
	}
	sType := model.ServiceType{
		ID:   int(item.Id),
		Name: item.Name,
	}
	c.JSON(http.StatusOK, sType)
}
