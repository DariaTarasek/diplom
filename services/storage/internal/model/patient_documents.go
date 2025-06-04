package model

import (
	"github.com/google/uuid"
	"time"
)

type (
	PatientDocument struct {
		ID          uuid.UUID  `db:"id"`
		PatientID   UserID     `db:"patient_id"`
		FileName    string     `db:"file_name"`
		Modality    string     `db:"modality"`
		StudyDate   *time.Time `db:"study_date"`
		Description string     `db:"description"`
		StoragePath string     `db:"storage_path"`
		PreviewPath string     `db:"preview_path"`
		CreatedAt   time.Time  `db:"created_at"`
	}
)
