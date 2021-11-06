package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/config"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/driver"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/repository"
)

type ReportServiceImpl struct {
	*config.App
	repository.ReportRepository
	repository.UserRepository
	PredictAPIURL string
}

func NewReportService(
	app *config.App,
	reportRepo repository.ReportRepository,
	userRepo repository.UserRepository,
	predictAPIURL string) ReportService {
	return &ReportServiceImpl{
		App:              app,
		ReportRepository: reportRepo,
		UserRepository:   userRepo,
		PredictAPIURL:    predictAPIURL,
	}
}

func (s *ReportServiceImpl) Create(
	ctx context.Context,
	report *entity.Report,
	image multipart.File,
	header *multipart.FileHeader) (*entity.Report, error) {
	format := strings.Split(header.Filename, ".")
	const op = "ReportServiceImpl.Create"
	if !allowedFileFormats(format[len(format)-1]) {
		return nil, api.NewSingleMessageException(
			api.EINVALID,
			op,
			fmt.Sprintf("%s format not supported", format),
			errors.New("file format not supported"),
		)
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", header.Filename)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"writer.CreateFormFile",
			err,
		)
	}
	_, err = io.Copy(part, image)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"io.Copy",
			err,
		)
	}
	writer.Close()
	timeoutCTX, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCTX, http.MethodPost, s.PredictAPIURL, body)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"http.NewRequestWithContext",
			err,
		)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		if uerr, ok := err.(*url.Error); ok && uerr.Timeout() {
			return nil, api.NewSingleMessageException(
				api.EUNAVAILABLE,
				op,
				"Timed out when trying to predict image. Please Try Again",
				err,
			)
		}
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"client.Do",
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &api.Exception{
			Op:  op,
			Err: errors.New("prediction service not returning 200 OK"),
		}
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"ioutil.ReadAll",
			err,
		)
	}
	predictResult := struct {
		Data *model.PredictResult `json:"data"`
	}{}
	if err := json.Unmarshal(resBody, &predictResult); err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"json.Unmarshal",
			err,
		)
	}
	report.ImageURL = predictResult.Data.ImageUrl
	report.Classes = predictResult.Data.Classes

	user := new(entity.User)
	driver.WithTransaction(s.App.DB, func(e driver.Executor) error {
		user, err = s.UserRepository.Get(ctx, s.App.DB, report.UserID)
		if err != nil {
			return err
		}
		report, err = s.ReportRepository.Create(ctx, s.App.DB, user.ID, report)
		if err != nil {
			return err
		}
		report.ReporterName = user.Name

		return nil
	})

	return report, nil
}

func (s *ReportServiceImpl) GetAll(ctx context.Context, pagination *model.Pagination) ([]*entity.Report, error) {
	return s.ReportRepository.GetAll(ctx, s.App.DB, pagination)
}

func (s *ReportServiceImpl) GetAllByUserID(ctx context.Context, userID int, pagination *model.Pagination) ([]*entity.Report, error) {
	return s.ReportRepository.GetAllByUserID(ctx, s.App.DB, userID, pagination)
}

func (s *ReportServiceImpl) Update(ctx context.Context, status string, reportID int) (*entity.Report, error) {
	return s.ReportRepository.Update(ctx, s.App.DB, status, reportID)
}

func allowedFileFormats(format string) bool {
	formats := []string{"jpg", "jpeg", "png"}
	for _, v := range formats {
		if v == format {
			return true
		}
	}

	return false
}
