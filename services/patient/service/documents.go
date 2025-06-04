package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/DariaTarasek/diplom/services/patient/model"
	authpb "github.com/DariaTarasek/diplom/services/patient/proto/auth"
	storagepb "github.com/DariaTarasek/diplom/services/patient/proto/storage"
	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	"google.golang.org/protobuf/types/known/timestamppb"
	"image"
	"image/color"
	"image/jpeg"
	"strconv"
	"strings"
	"time"
)

func (s *PatientService) UploadTest(ctx context.Context, input model.UploadTestInput) (uuid.UUID, error) {
	// 1. Проверка, что это DICOM
	if !isDICOM(input.FileContent) {
		return uuid.Nil, errors.New("файл не является файлом DICOM")
	}

	if !hasExtension(input.FileName) {
		input.FileName += ".dcm"
	}

	// 2. Получение ID пациента
	authResp, err := s.AuthClient.Client.GetPatient(ctx, &authpb.GetPatientRequest{
		Token: input.Token,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("не удалось получить идентификатор пациента: %w", err)
	}
	patientID := int(authResp.Patient.UserId)

	// 3. Извлекаем метаданные из DICOM
	meta, err := extractDICOMMetadata(input.FileContent)
	if err != nil {
		return uuid.Nil, fmt.Errorf("не удалось извлечь информацию из файла DICOM: %w", err)
	}

	// 4. Генерируем превью
	preview, err := generatePreviewJPEG(input.FileContent)
	if err != nil {
		// Можно логировать, но не критично
		preview = nil
	}

	patID := strconv.Itoa(patientID)
	// 5. Вызываем storage.SaveDocument
	storageResp, err := s.StorageClient.Client.SaveDocument(ctx, &storagepb.SaveDocumentRequest{
		PatientId:       patID,
		FileName:        input.FileName,
		FileContent:     input.FileContent,
		Modality:        meta.Modality,
		StudyDate:       timestamppb.New(derefTime(meta.StudyDate)),
		Description:     input.Description,
		PreviewJpeg:     preview,
		PreviewFileName: input.FileName[:len(input.FileName)-4] + ".jpg",
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("storage error: %w", err)
	}

	return uuid.MustParse(storageResp.DocumentId), nil
}

func (s *PatientService) GetDocumentsInfo(ctx context.Context, token string) ([]model.DocumentInfo, error) {
	userID, err := s.AuthClient.Client.GetPatient(ctx, &authpb.GetPatientRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить идентификатор пользователя: %w", err)
	}

	uid := strconv.Itoa(int(userID.Patient.UserId))
	resp, err := s.StorageClient.Client.GetDocumentsByPatientID(ctx, &storagepb.GetDocumentsRequest{PatientId: uid})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о документах: %w", err)
	}

	docsInfo := make([]model.DocumentInfo, 0, len(resp.Documents))
	for _, item := range resp.Documents {
		docInfo := model.DocumentInfo{
			ID:          item.Id,
			FileName:    item.FileName,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.AsTime().Format("02.01.2006"),
		}
		docsInfo = append(docsInfo, docInfo)
	}
	return docsInfo, nil
}

func (s *PatientService) DownloadDocument(ctx context.Context, documentID string) (model.DocumentFile, error) {

	// Запрашиваем документ по ID через storage-сервис
	resp, err := s.StorageClient.Client.DownloadDocument(ctx, &storagepb.DownloadDocumentRequest{DocumentId: documentID})
	if err != nil {
		return model.DocumentFile{}, fmt.Errorf("не удалось скачать документ: %w", err)
	}

	// Собираем и возвращаем структуру модели
	return model.DocumentFile{
		FileName:    resp.FileName,
		FileContent: resp.FileContent,
	}, nil
}

func isDICOM(data []byte) bool {
	if len(data) < 132 {
		return false
	}
	return string(data[128:132]) == "DICM"
}

func extractDICOMMetadata(data []byte) (*model.PatientDocument, error) {
	dataset, err := dicom.ParseUntilEOF(bytes.NewReader(data), nil)
	if err != nil {
		return nil, err
	}

	var modality, description string
	var studyDate *time.Time

	if elem, err := dataset.FindElementByTag(tag.Modality); err == nil {
		modality = elem.String()
	}

	if elem, err := dataset.FindElementByTag(tag.SeriesDescription); err == nil {
		description = elem.String()
	}

	if elem, err := dataset.FindElementByTag(tag.StudyDate); err == nil {
		dateStr := elem.String()
		if t, err := time.Parse("20060102", dateStr); err == nil {
			studyDate = &t
		}
	}

	return &model.PatientDocument{
		Modality:    modality,
		StudyDate:   studyDate,
		Description: description,
	}, nil
}

func generatePreviewJPEG(fileContent []byte) ([]byte, error) {
	// Parsing DICOM from bytes
	dataset, err := dicom.Parse(bytes.NewReader(fileContent), int64(len(fileContent)), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге DICOM: %v", err)
	}

	// Getting pixel data
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске PixelData: %v", err)
	}

	// Getting pixel data information
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)

	// Checking for frames
	if len(pixelDataInfo.Frames) == 0 {
		return nil, fmt.Errorf("нет фреймов в DICOM файле")
	}

	// Using the first frame for preview
	frame := pixelDataInfo.Frames[0]

	// Getting image from frame
	img, err := frame.GetImage()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении изображения из фрейма: %v", err)
	}

	// Normalizing image intensity
	normalizedImg := normalizeImage(img)

	// Buffer for encoding JPG
	var buf bytes.Buffer

	// Encoding image to JPG
	err = jpeg.Encode(&buf, normalizedImg, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, fmt.Errorf("ошибка при кодировании JPG: %v", err)
	}

	return buf.Bytes(), nil
}

// normalizeImage нормализует интенсивность изображения для улучшения контрастности
func normalizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	normalized := image.NewGray(bounds)

	// Finding min and max intensity
	minVal, maxVal := findMinMaxIntensity(img)

	// Avoiding division by zero
	if maxVal == minVal {
		maxVal = minVal + 1
	}

	// Linear normalization to 0-255
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Assuming grayscale (common for X-ray)
			grayVal := getGrayValue(img, x, y)
			// Scaling to 0-255
			normalizedVal := uint8(255 * (float64(grayVal) - float64(minVal)) / (float64(maxVal) - float64(minVal)))
			normalized.SetGray(x, y, color.Gray{Y: normalizedVal})
		}
	}

	return normalized
}

// findMinMaxIntensity находит минимальную и максимальную интенсивность в изображении
func findMinMaxIntensity(img image.Image) (min, max uint16) {
	bounds := img.Bounds()
	min, max = 65535, 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayVal := getGrayValue(img, x, y)
			if grayVal < min {
				min = grayVal
			}
			if grayVal > max {
				max = grayVal
			}
		}
	}
	return
}

// getGrayValue извлекает значение яркости пикселя (для 16-bit grayscale)
func getGrayValue(img image.Image, x, y int) uint16 {
	switch img := img.(type) {
	case *image.Gray:
		return uint16(img.GrayAt(x, y).Y)
	case *image.Gray16:
		return img.Gray16At(x, y).Y
	default:
		// Fallback for other image types
		r, g, b, _ := img.At(x, y).RGBA()
		// Converting to grayscale using luminosity method
		return uint16(0.299*float64(r)+0.587*float64(g)+0.114*float64(b)) >> 8
	}
}

func hasExtension(fileName string) bool {
	dot := strings.LastIndex(fileName, ".")
	return dot != -1 && dot != len(fileName)-1
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
