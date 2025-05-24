package model

type (
	ServiceTypeID int
	ServiceID     int
	ServiceType   struct {
		ID   ServiceTypeID `db:"id"`
		Name string        `db:"name"`
	}
	Service struct {
		ID       ServiceID     `db:"id"`
		Name     string        `db:"name"`
		Price    *int          `db:"price"`
		Category ServiceTypeID `db:"type"`
	}
)
