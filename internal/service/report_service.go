package service

import (
	"mime/multipart"

	"github.com/muktiarafi/rodavis-api/internal/entity"
)

type ReportService interface {
	Create(
		report *entity.Report,
		image multipart.File,
		header *multipart.FileHeader,
	) (*entity.Report, error)
	GetAll(limit, lastseenID uint64) ([]*entity.Report, error)
	GetAllByUserID(userID int, limit, lastseenID uint64) ([]*entity.Report, error)
	Update(status string, reportID int) (*entity.Report, error)
}
