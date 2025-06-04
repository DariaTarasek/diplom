package patient

import (
	"fmt"
	"github.com/DariaTarasek/diplom/services/api-gateway/model"
	patientpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/patient"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

func (h *PatientHandler) UploadTest(c *gin.Context) {
	// 1. Получаем токен из куки
	token, err := c.Cookie("access_token")
	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "необходима авторизация"})
		return
	}

	// 2. Читаем multipart/form-data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден в запросе " + err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть файл " + err.Error()})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать файл " + err.Error()})
		return
	}

	description := c.PostForm("description")

	// 3. Вызов gRPC метода UploadTest
	resp, err := h.PatientClient.Client.UploadTest(c.Request.Context(), &patientpb.UploadTestRequest{
		Token:       token,
		FileName:    fileHeader.Filename,
		FileContent: fileBytes,
		Description: description,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка " + err.Error()})
		}
		return
	}

	// 4. Успех
	c.JSON(http.StatusOK, gin.H{"document_id": resp.DocumentId})
}

func (h *PatientHandler) getDocuments(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "необходима авторизация"})
		return
	}

	resp, err := h.PatientClient.Client.GetDocumentsByPatientID(c.Request.Context(), &patientpb.GetDocumentsRequest{Token: token})
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

func (h *PatientHandler) DownloadDocument(c *gin.Context) {
	// Получаем токен и ID документа из запроса
	documentID := c.Param("id") // например: /documents/:id

	// Получаем файл через gRPC
	doc, err := h.PatientClient.Client.DownloadDocument(c.Request.Context(), &patientpb.DownloadDocumentRequest{DocumentId: documentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отдаем файл пользователю
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", doc.FileName))
	c.Data(http.StatusOK, "application/octet-stream", doc.FileContent)
}
