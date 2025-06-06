package doctor

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getDoctorProfile godoc
// @Summary Получение профиля врача
// @Tags Врач
// @Success 200 {object} model.Doctor
// @Failure 401 {object} gin.H "Необходима авторизация"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H
// @Router /doctor/me [get]
func (h *DoctorHandler) getDoctorProfile(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "необходима авторизация"})
		return
	}

	resp, err := h.AuthClient.Client.GetDoctorProfile(c.Request.Context(), &authpb.GetProfileRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	exp := int(resp.Doctor.Experience)
	doc := model.Doctor{
		ID:          int(resp.Doctor.UserId),
		FirstName:   resp.Doctor.FirstName,
		SecondName:  resp.Doctor.SecondName,
		Surname:     &resp.Doctor.Surname,
		PhoneNumber: &resp.Doctor.PhoneNumber,
		Email:       resp.Doctor.Email,
		Education:   &resp.Doctor.Education,
		Experience:  &exp,
		Gender:      resp.Doctor.Gender,
	}

	c.JSON(http.StatusOK, doc)
}
