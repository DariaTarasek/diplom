package doctor

import (
	"fmt"
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	doctorpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/doctor"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetPatientAllergiesChronics godoc
// @Summary Получение аллергий и хронических заболеваний пациента
// @Tags Врач
// @Param id path int true "ID пациента"
// @Success 200 {array} model.AllergiesChronics
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Router /api/patient-notes/{id} [get]
func (h *DoctorHandler) GetPatientAllergiesChronics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	notesResp, err := h.DoctorClient.Client.GetPatientAllergiesChronics(c.Request.Context(), &doctorpb.GetByIdRequest{Id: int32(id)})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	var notes []model.AllergiesChronics
	for _, item := range notesResp.PatientAllergiesChronics {
		note := model.AllergiesChronics{
			ID:        int(item.Id),
			PatientID: int(item.PatientId),
			Type:      item.Type,
			Title:     item.Title,
		}
		notes = append(notes, note)
	}
	c.JSON(http.StatusOK, notes)
}

// GetAppointmentByID godoc
// @Summary Получение информации о приёме по ID
// @Tags Врач
// @Param id path int true "ID приёма"
// @Success 200 {object} model.Appointment
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Router /api/appointments/{id} [get]
func (h *DoctorHandler) GetAppointmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	apptResp, err := h.DoctorClient.Client.GetAppointmentByID(c.Request.Context(), &doctorpb.GetByIdRequest{Id: int32(id)})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	patientID := model.UserID(apptResp.Appt.PatientId)
	appt := model.Appointment{
		ID:                 model.AppointmentID(apptResp.Appt.Id),
		DoctorID:           model.UserID(apptResp.Appt.DoctorId),
		PatientID:          &patientID,
		Date:               apptResp.Appt.Date.AsTime().Format("02.01.06"),
		Time:               apptResp.Appt.Time.AsTime().Format("15:04"),
		PatientSecondName:  apptResp.Appt.SecondName,
		PatientFirstName:   apptResp.Appt.FirstName,
		PatientSurname:     &apptResp.Appt.Surname,
		PatientBirthDate:   apptResp.Appt.BirthDate.AsTime().Format("02.01.2006"),
		PatientGender:      apptResp.Appt.Gender,
		PatientPhoneNumber: apptResp.Appt.PhoneNumber,
		Status:             apptResp.Appt.Status,
		CreatedAt:          apptResp.Appt.CreatedAt.AsTime().Format("2006.01.02 15:04"),
		UpdatedAt:          apptResp.Appt.UpdatedAt.AsTime().Format("2006.01.02 15:04"),
	}

	c.JSON(http.StatusOK, appt)
}

// GetPatientVisits godoc
// @Summary Получение списка посещений пациента
// @Tags Врач
// @Param id path int true "ID пациента"
// @Success 200 {array} model.Visit
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Router /api/patient-history/{id} [get]
func (h *DoctorHandler) GetPatientVisits(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	visitsResp, err := h.DoctorClient.Client.GetPatientVisits(c.Request.Context(), &doctorpb.GetByIdRequest{Id: int32(id)})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	var visits []model.Visit
	for _, item := range visitsResp.Visits {
		visit := model.Visit{
			ID:            int(item.Id),
			AppointmentID: int(item.ApptId),
			PatientID:     int(item.PatientId),
			Doctor:        item.Doctor,
			Complaints:    item.Complaints,
			Treatment:     item.Treatment,
			CreatedAt:     item.CreatedAt,
			Diagnoses:     DerefDiagnoses(item.Diagnoses),
		}
		visits = append(visits, visit)
	}
	fmt.Println(visits)
	c.JSON(http.StatusOK, visits)
}

func DerefDiagnoses(src []*doctorpb.Diagnose) []model.Diagnose {
	result := make([]model.Diagnose, 0, len(src))
	for _, d := range src {
		if d != nil {
			result = append(result, model.Diagnose{
				ICDCode: d.IcdCode,
				Notes:   d.Notes,
			})
		}
	}
	return result
}

// AddPatientAllergiesChronics godoc
// @Summary Добавление аллергий и хронических заболеваний пациенту
// @Tags Врач
// @Param id path int true "ID пациента"
// @Param body body []model.AllergiesChronics true "Список аллергий и хронических заболеваний"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H "Ошибка сервера"
// @Router /api/patient-notes/{id} [post]
func (h *DoctorHandler) AddPatientAllergiesChronics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var notes []model.AllergiesChronics
	if err := c.ShouldBindJSON(&notes); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	var notesReq []*doctorpb.PatientAllergiesChronics
	for _, item := range notes {
		note := &doctorpb.PatientAllergiesChronics{
			PatientId: int32(id),
			Type:      item.Type,
			Title:     item.Title,
		}
		notesReq = append(notesReq, note)
	}
	AddNotesRequest := &doctorpb.AddPatientAllergiesChronicsRequest{Notes: notesReq}
	_, err = h.DoctorClient.Client.AddPatientAllergiesChronics(c.Request.Context(), AddNotesRequest)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// AddConsultation godoc
// @Summary Добавление консультации
// @Tags Врач
// @Param body body model.VisitSaveRequest true "Данные консультации"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H "Некорректный ввод"
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Router /api/visits [post]
func (h *DoctorHandler) AddConsultation(c *gin.Context) {
	var visit model.VisitSaveRequest
	if err := c.ShouldBindJSON(&visit); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var materialsReq []*doctorpb.AddVisitMaterials
	for _, item := range visit.Materials {
		material := &doctorpb.AddVisitMaterials{
			MaterialId: int32(item.ID),
			Amount:     int32(item.Quantity),
		}
		materialsReq = append(materialsReq, material)
	}

	var servicesReq []*doctorpb.AddVisitServices
	for _, item := range visit.Services {
		service := &doctorpb.AddVisitServices{
			ServiceId: int32(item.ID),
			Amount:    int32(item.Quantity),
		}
		servicesReq = append(servicesReq, service)
	}

	var diagnosesReq []*doctorpb.VisitDiagnose
	for _, item := range visit.ICDCodes {
		diagnose := &doctorpb.VisitDiagnose{
			IcdCodeId: int32(item.CodeID),
			Note:      item.Comment,
		}
		diagnosesReq = append(diagnosesReq, diagnose)
	}

	_, err := h.DoctorClient.Client.AddConsultation(c.Request.Context(), &doctorpb.AddConsultationRequest{
		AppointmentId: int32(visit.AppointmentID),
		PatientId:     int32(visit.PatientID),
		DoctorId:      int32(visit.DoctorID),
		Complaints:    visit.Complaints,
		Treatment:     visit.Treatment,
		Diagnoses:     diagnosesReq,
		Services:      servicesReq,
		Materials:     materialsReq,
	})
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// getPatientDocs godoc
// @Summary Получение документов пациента
// @Tags Врач
// @Param id path int true "ID пациента"
// @Success 200 {array} model.DocumentInfo
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H
// @Router /admin/patient/{id}/documents [get]
func (h *DoctorHandler) getPatientDocs(c *gin.Context) {
	patientId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.DoctorClient.Client.GetDocumentsByPatientID(c.Request.Context(), &doctorpb.GetDocumentsRequest{PatientID: int32(patientId)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	docs := make([]model.DocumentInfo, 0, len(resp.Documents))
	for _, item := range resp.Documents {
		doc := model.DocumentInfo{
			ID:          item.Id,
			FileName:    item.FileName,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		}
		docs = append(docs, doc)
	}
	c.JSON(http.StatusOK, docs)
}

// DownloadDocument godoc
// @Summary Скачивание документа пациента
// @Tags Врач
// @Param id path string true "ID документа"
// @Success 200 {file} file
// @Failure 403 {object} gin.H "Недостаточно прав"
// @Failure 500 {object} gin.H
// @Router /api/doctor/consultation/patient-tests/{id} [get]
func (h *DoctorHandler) DownloadDocument(c *gin.Context) {
	// Получаем токен и ID документа из запроса
	documentID := c.Param("id") // например: /documents/:id

	// Получаем файл через gRPC
	doc, err := h.DoctorClient.Client.DownloadDocument(c.Request.Context(), &doctorpb.DownloadDocumentRequest{DocumentId: documentID})
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отдаем файл пользователю
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", doc.FileName))
	c.Data(http.StatusOK, "application/octet-stream", doc.FileContent)
}
