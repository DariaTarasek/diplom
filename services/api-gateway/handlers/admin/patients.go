package admin

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	adminpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/admin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetPatients(c *gin.Context) {
	items, err := h.AdminClient.Client.GetPatients(c.Request.Context(), &adminpb.EmptyRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var patients []model.PatientWithoutPassword
	for _, item := range items.Patients {
		patient := model.PatientWithoutPassword{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     &item.Surname,
			PhoneNumber: &item.PhoneNumber,
			Email:       &item.Email,
			Gender:      item.Gender,
			BirthDate:   item.BirthDate,
		}
		patients = append(patients, patient)
	}
	c.JSON(http.StatusOK, patients)
}

func (h *Handler) UpdatePatient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var patientReq model.Patient
	if err := c.ShouldBindJSON(&patientReq); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	UpdatePatientRequest := &adminpb.UpdatePatientRequest{
		UserId:      int32(id),
		FirstName:   patientReq.FirstName,
		Surname:     *patientReq.Surname,
		SecondName:  patientReq.SecondName,
		Email:       *patientReq.Email,
		BirthDate:   patientReq.BirthDate,
		PhoneNumber: patientReq.PhoneNumber,
		Gender:      patientReq.Gender,
	}

	_, err = h.AdminClient.Client.UpdatePatient(c.Request.Context(), UpdatePatientRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
