package model

type VisitPayment struct {
	VisitID int
	Price   int
	Status  string
}

type UnconfirmedVisitPayment struct {
	VisitID   int
	Doctor    string
	Patient   string
	CreatedAt string
	Price     int
}

type VisitMaterialsServices struct {
	ID       int
	VisitID  int
	Item     string
	Quantity int
}
