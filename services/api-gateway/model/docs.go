package model

type (
	DocumentInfo struct {
		ID          string `json:"id"`
		FileName    string `json:"file_name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
	}
	DocumentFile struct {
		FileName    string
		FileContent []byte
	}
)
