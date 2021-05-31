package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/muktiarafi/rodavis-api/internal/api"
	"github.com/muktiarafi/rodavis-api/internal/driver"
	"github.com/muktiarafi/rodavis-api/internal/entity"
)

type ReportRepositoryImpl struct {
	*driver.DB
}

func NewReportRepository(db *driver.DB) ReportRepository {
	return &ReportRepositoryImpl{
		DB: db,
	}
}

func (r *ReportRepositoryImpl) Create(userID int, report *entity.Report) (*entity.Report, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	stmt := `INSERT INTO reports (image_url, classes, note, address, lat, lng, user_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, status, image_url, classes, note, address, lat, lng, date_reported`

	var cls pgtype.EnumArray
	if err := r.SQL.QueryRowContext(
		ctx,
		stmt,
		report.ImageURL,
		report.Classes,
		report.Note,
		report.Address,
		report.Location.Lat,
		report.Location.Lng,
		report.UserID,
	).Scan(
		&report.ID,
		&report.Status,
		&report.ImageURL,
		&cls,
		&report.Note,
		&report.Address,
		&report.Location.Lat,
		&report.Location.Lng,
		&report.DateReported,
	); err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			"ReportRepositoryImpl.Create",
			"r.SQL.QueryRowContext",
			err,
		)
	}
	classes := make([]string, len(cls.Elements))
	for k, v := range cls.Elements {
		classes[k] = v.String
	}
	report.Classes = classes

	return report, nil
}

func (r *ReportRepositoryImpl) GetAll(limit, lastseenID uint64) ([]*entity.Report, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	const op = "ReportRepositoryImpl.GetAll"
	queryBuilder := squirrel.
		Select("r.id", "name", "status", "image_url", "note", "address", "lat", "lng", "date_reported").
		From("users AS u").Join("reports AS r ON u.id = r.user_id").PlaceholderFormat(squirrel.Dollar).OrderBy("r.id DESC")

	if limit > 0 {
		queryBuilder = queryBuilder.Limit(limit)
	}

	if lastseenID > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Lt{
			"r.id": lastseenID,
		})
	}

	stmt, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"queryBuilder.ToSql",
			err,
		)
	}

	rows, err := r.SQL.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"r.SQL.QueryContext",
			err,
		)
	}

	defer rows.Close()
	reports := []*entity.Report{}
	for rows.Next() {
		report := new(entity.Report)
		location := new(entity.Location)
		var cls pgtype.EnumArray
		if err := rows.Scan(
			&report.ID,
			&report.ReporterName,
			&report.Status,
			&report.ImageURL,
			&report.Note,
			&report.Address,
			&location.Lat,
			&location.Lng,
			&report.Address,
		); err != nil {
			return nil, api.NewExceptionWithSourceLocation(
				op,
				"rows.Scan",
				err,
			)
		}
		classes := make([]string, len(cls.Elements))
		for k, v := range cls.Elements {
			classes[k] = v.String
		}
		report.Classes = classes
		report.Location = location
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportRepositoryImpl) GetAllByUserID(userID int, limit, lastseenID uint64) ([]*entity.Report, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	const op = "ReportRepositoryImpl.GetAll"
	queryBuilder := squirrel.
		Select("r.id", "name", "status", "image_url", "note", "address", "lat", "lng", "date_reported").
		From("users AS u").Join("reports AS r ON u.id = r.user_id").PlaceholderFormat(squirrel.Dollar).OrderBy("r.id DESC").
		Where(squirrel.Eq{
			"r.user_id": userID,
		})

	if limit > 0 {
		queryBuilder = queryBuilder.Limit(limit)
	}

	if lastseenID > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Lt{
			"r.id": lastseenID,
		})
	}

	stmt, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"queryBuilder.ToSql",
			err,
		)
	}

	rows, err := r.SQL.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"r.SQL.QueryContext",
			err,
		)
	}

	defer rows.Close()
	reports := []*entity.Report{}
	for rows.Next() {
		report := new(entity.Report)
		location := new(entity.Location)
		var cls pgtype.EnumArray
		if err := rows.Scan(
			&report.ID,
			&report.ReporterName,
			&report.Status,
			&report.ImageURL,
			&report.Note,
			&report.Address,
			&location.Lat,
			&location.Lng,
			&report.Address,
		); err != nil {
			return nil, api.NewExceptionWithSourceLocation(
				op,
				"rows.Scan",
				err,
			)
		}
		classes := make([]string, len(cls.Elements))
		for k, v := range cls.Elements {
			classes[k] = v.String
		}
		report.Classes = classes
		report.Location = location
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportRepositoryImpl) Update(status string, reportID int) (*entity.Report, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	stmt := `UPDATE reports
	SET status = $1
	WHERE id = $2
	RETURNING id, status, image_url, classes, note, address, lat, lng, date_reported`

	report := new(entity.Report)
	location := new(entity.Location)
	var cls pgtype.EnumArray
	if err := r.SQL.QueryRowContext(ctx, stmt, status, reportID).Scan(
		&report.ID,
		&report.Status,
		&report.ImageURL,
		&cls,
		&report.Note,
		&report.Address,
		&location.Lat,
		&location.Lng,
		&report.DateReported,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, api.NewSingleMessageException(
				api.ENOTFOUND,
				"ReportRepositoryImpl.Update",
				"Report Not Found",
				err,
			)
		}
		return nil, api.NewExceptionWithSourceLocation(
			"ReportRepositoryImpl.Update",
			"r.SQL.QueryRowContext",
			err,
		)
	}
	classes := make([]string, len(cls.Elements))
	for k, v := range cls.Elements {
		classes[k] = v.String
	}
	report.Classes = classes
	report.Location = location

	return report, nil
}