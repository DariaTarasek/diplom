package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getAdminProfile(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "необходима авторизация"})
		return
	}

	resp, err := h.AuthClient.Client.GetAdminProfile(c.Request.Context(), &authpb.GetProfileRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	admin := model.AdminWithRole{
		ID:          int(resp.Admin.UserId),
		FirstName:   resp.Admin.FirstName,
		SecondName:  resp.Admin.SecondName,
		Surname:     resp.Admin.Surname,
		PhoneNumber: resp.Admin.PhoneNumber,
		Email:       resp.Admin.Email,
		Gender:      resp.Admin.Gender,
		Role:        resp.Admin.Role,
	}

	c.JSON(http.StatusOK, admin)
}
