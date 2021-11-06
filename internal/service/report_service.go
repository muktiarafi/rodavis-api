package service

import (
	"context"
	"mime/multipart"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

type ReportService interface {
	Create(
		ctx context.Context,
		report *entity.Report,
		image multipart.File,
		header *multipart.FileHeader,
	) (*entity.Report, error)
	GetAll(ctx context.Context, pagination *model.Pagination) ([]*entity.Report, error)
	GetAllByUserID(ctx context.Context, userID int, pagination *model.Pagination) ([]*entity.Report, error)
	Update(ctx context.Context, status string, reportID int) (*entity.Report, error)
}
