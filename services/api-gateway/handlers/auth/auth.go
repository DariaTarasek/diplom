package auth

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type roleResponse struct {
	Role string `json:"role"`
}

func (h *Handler) authorize(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.AuthClient.Client.Auth(c.Request.Context(), &authpb.AuthRequest{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	setTokenCookie(c, resp.Token)

	c.JSON(http.StatusOK, roleResponse{Role: resp.Role})
}

func setTokenCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
