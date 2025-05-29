package model

type ServiceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	ServiceTypeId int    `json:"category_id"`
}
