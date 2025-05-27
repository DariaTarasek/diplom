package model

const (
	SuperAdminRole = 1
	AdminRole      = 2
	DoctorRole     = 3
	PatientRole    = 4
)

type (
	PermissionID int
	RoleID       int
	Role         struct {
		ID   RoleID
		Name string
	}
	UserRole struct {
		UserID UserID
		RoleID RoleID
	}
)
