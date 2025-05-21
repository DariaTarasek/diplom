package model

type (
	RoleID       int
	PermissionID int
	Role         struct {
		ID   RoleID `db:"id"`
		Name string `db:"name"`
	}
	Permission struct {
		ID   PermissionID `db:"id"`
		Name string       `db:"name"`
	}
	RolePermissions struct {
		RoleID       RoleID       `db:"role_id"`
		PermissionID PermissionID `db:"permission_id"`
	}
	UserRole struct {
		UserID UserID `db:"user_id"`
		RoleID RoleID `db:"role_id"`
	}
)
