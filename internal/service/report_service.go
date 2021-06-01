package service

import (
	"context"
	"mime/multipart"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
)

type ReportService interface {
	Create(
		ctx context.Context,
		report *entity.Report,
		image multipart.File,
		header *multipart.FileHeader,
	) (*entity.Report, error)
	GetAll(ctx context.Context, limit, lastseenID uint64) ([]*entity.Report, error)
	GetAllByUserID(ctx context.Context, userID int, limit, lastseenID uint64) ([]*entity.Report, error)
	Update(ctx context.Context, status string, reportID int) (*entity.Report, error)
}
