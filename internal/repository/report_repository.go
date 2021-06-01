package repository

import (
	"context"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
)

type ReportRepository interface {
	Create(ctx context.Context, userID int, report *entity.Report) (*entity.Report, error)
	GetAll(ctx context.Context, limit, lastseenID uint64) ([]*entity.Report, error)
	GetAllByUserID(ctx context.Context, userID int, limit, lastseenID uint64) ([]*entity.Report, error)
	Update(ctx context.Context, status string, reportID int) (*entity.Report, error)
}
