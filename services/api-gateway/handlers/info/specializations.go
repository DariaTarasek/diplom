package info

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	storagepb "github.com/DariaTarasek/diplom/services/api-gateway/proto/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *InfoHandler) GetAllSpecs(c *gin.Context) {
	items, err := h.store.Client.GetAllSpecs(c.Request.Context(), &storagepb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var specs []model.Specialization
	for _, item := range items.Specs {
		spec := model.Specialization{
			ID:   int(item.Id),
			Name: item.Name,
		}
		specs = append(specs, spec)
	}
	c.JSON(http.StatusOK, specs)
}
