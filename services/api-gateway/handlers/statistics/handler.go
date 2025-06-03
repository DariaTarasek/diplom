package statistics

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	StatisticsClient *clients.StatisticsClient
}

func NewHandler(statClient *clients.StatisticsClient) *Handler {
	return &Handler{
		StatisticsClient: statClient,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/statistics", h.GetAllStats)
}
