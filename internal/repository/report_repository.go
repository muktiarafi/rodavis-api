package repository

import "gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"

type ReportRepository interface {
	Create(userID int, report *entity.Report) (*entity.Report, error)
	GetAll(limit, lastseenID uint64) ([]*entity.Report, error)
	GetAllByUserID(userID int, limit, lastseenID uint64) ([]*entity.Report, error)
	Update(status string, reportID int) (*entity.Report, error)
}
