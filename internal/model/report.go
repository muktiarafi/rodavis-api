package model

type CreateReportDTO struct {
	Lat     string `validate:"required"`
	Lng     string `validate:"required"`
	Image   int64  `validate:"required"`
	Address string `validate:"required,min=4"`
}

type UpdateReportDTO struct {
	Status string `json:"status" validate:"oneof='Reported' 'Under Repair' 'Completed' 'Rejected'"`
}
