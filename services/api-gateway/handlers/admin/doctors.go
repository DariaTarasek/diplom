package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetSpecs godoc
// @Summary Получить список специализаций
// @Tags Администратор
// @Description Возвращает все доступные специализации
// @Produce json
// @Success 200 {array} model.Spec
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/specialties [get]
func (h *Handler) GetSpecs(c *gin.Context) {
	items, err := h.AdminClient.Client.GetSpecs(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var specs []model.Spec
	for _, item := range items.Specs {
		spec := model.Spec{
			ID:   int(item.Id),
			Name: item.Name,
		}
		specs = append(specs, spec)
	}
	c.JSON(http.StatusOK, specs)
}

// GetDoctors godoc
// @Summary Получить список докторов
// @Tags Администратор
// @Description Возвращает всех докторов с их специализациями
// @Produce json
// @Success 200 {array} model.DoctorWithSpecs
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/staff-doctors [get]
func (h *Handler) GetDoctors(c *gin.Context) {
	items, err := h.AdminClient.Client.GetDoctors(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var doctors []model.DoctorWithSpecs
	for _, item := range items.Doctors {
		var specs []int
		for _, specItem := range item.Specs {
			intSpecItem := int(specItem)
			specs = append(specs, intSpecItem)
		}
		doctor := model.DoctorWithSpecs{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Education:   item.Education,
			Experience:  int(item.Experience),
			Gender:      item.Gender,
			Specs:       specs,
		}
		doctors = append(doctors, doctor)
	}
	c.JSON(http.StatusOK, doctors)
}

// UpdateDoctor godoc
// @Summary Обновить данные доктора
// @Tags Администратор
// @Description Обновляет информацию о докторе и его специализациях
// @Accept json
// @Produce json
// @Param doctor body model.DoctorWithSpecs true "Данные доктора"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Доступ запрещён"
// @Failure 500 {object} gin.H "Внутренняя ошибка сервера"
// @Router /api/save-doctor [put]
func (h *Handler) UpdateDoctor(c *gin.Context) {
	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	log.Println(err.Error())
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
	//	return
	//}
	var doctorReq model.DoctorWithSpecs
	if err := c.ShouldBindJSON(&doctorReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var specsInt32 []int32
	for _, item := range doctorReq.Specs {
		specsInt32 = append(specsInt32, int32(item))
	}

	UpdateDoctorRequest := &adminpb.UpdateDoctorRequest{
		UserId:      int32(doctorReq.ID),
		FirstName:   doctorReq.FirstName,
		SecondName:  doctorReq.SecondName,
		Surname:     doctorReq.Surname,
		PhoneNumber: doctorReq.PhoneNumber,
		Email:       doctorReq.Email,
		Education:   doctorReq.Education,
		Experience:  int32(doctorReq.Experience),
		Gender:      doctorReq.Gender,
		Specs:       specsInt32,
	}

	_, err := h.AdminClient.Client.UpdateDoctor(c.Request.Context(), UpdateDoctorRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
