package model

type (
	ICDCodeID int
	ICD       struct {
		ID          ICDCodeID `db:"id"`
		Code        string    `db:"code"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
	}
	Diagnose struct {
		ID            int       `db:"id"`
		VisitID       VisitID   `db:"visit_id"`
		ICDCodeID     ICDCodeID `db:"icd_code_id"`
		DiagnosisNote string    `db:"diagnosis_note"`
	}
)
