package repository

import (
	"context"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/driver"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
)

type ReportRepository interface {
	Create(ctx context.Context, e driver.Executor, userID int, report *entity.Report) (*entity.Report, error)
	GetAll(ctx context.Context, e driver.Executor, pagination *model.Pagination) ([]*entity.Report, error)
	GetAllByUserID(ctx context.Context, e driver.Executor, userID int, pagination *model.Pagination) ([]*entity.Report, error)
	Update(ctx context.Context, e driver.Executor, status string, reportID int) (*entity.Report, error)
}
