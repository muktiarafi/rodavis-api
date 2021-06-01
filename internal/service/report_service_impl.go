package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/repository"
)

type ReportServiceImpl struct {
	repository.ReportRepository
	repository.UserRepository
	PredictAPIURL string
}

func NewReportService(
	reportRepo repository.ReportRepository,
	userRepo repository.UserRepository,
	predictAPIURL string) ReportService {
	return &ReportServiceImpl{
		ReportRepository: reportRepo,
		UserRepository:   userRepo,
		PredictAPIURL:    predictAPIURL,
	}
}

func (s *ReportServiceImpl) Create(
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
	res, err := http.Post(s.PredictAPIURL, writer.FormDataContentType(), body)
	if err != nil {
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"io.Copy",
			err,
		)
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

	user, err := s.UserRepository.Get(report.UserID)
	if err != nil {
		return nil, err
	}
	report, err = s.ReportRepository.Create(user.ID, report)
	if err != nil {
		return nil, err
	}
	report.ReporterName = user.Name

	return report, nil
}

func (s *ReportServiceImpl) GetAll(limit, lastseenID uint64) ([]*entity.Report, error) {
	return s.ReportRepository.GetAll(limit, lastseenID)
}

func (s *ReportServiceImpl) GetAllByUserID(userID int, limit, lastseenID uint64) ([]*entity.Report, error) {
	return s.ReportRepository.GetAllByUserID(userID, limit, lastseenID)
}

func (s *ReportServiceImpl) Update(status string, reportID int) (*entity.Report, error) {
	return s.ReportRepository.Update(status, reportID)
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
