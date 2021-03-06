package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/middleware"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/model"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/service"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/validation"
)

type ReportHandler struct {
	*validation.Validator
	service.ReportService
}

func NewReportHandler(
	val *validation.Validator,
	reportSRV service.ReportService,
) *ReportHandler {
	return &ReportHandler{
		Validator:     val,
		ReportService: reportSRV,
	}
}

func (h *ReportHandler) Route(mux *chi.Mux) {
	mux.Route("/api/reports", func(r chi.Router) {
		r.With(middleware.RequireAuth).Post("/", h.NewReport)
		r.Get("/", h.GetAllReport)
		r.With(middleware.RequireAuth).Get("/history", h.GetAllUserReport)
		r.With(middleware.RequireAuth).Put("/{reportID}", h.UpdateReport)
	})
}

func (h *ReportHandler) NewReport(w http.ResponseWriter, r *http.Request) {
	const op = "ReportHandler.NewReport"
	userPayload, err := api.UserPayloadFromContext(op, r)
	if err != nil {
		api.SendError(w, err)
		return
	}

	r.ParseMultipartForm(10 << 20)

	latStr := r.FormValue("lat")
	lngStr := r.FormValue("lng")
	address := r.FormValue("address")
	note := r.FormValue("note")

	image, header, err := r.FormFile("image")
	if err != nil {
		exc := api.NewSingleMessageException(
			api.EINVALID,
			op,
			"Image is Required",
			err,
		)
		api.SendError(w, exc)
		return
	}

	createReportDTO := &model.CreateReportDTO{
		Lat:     latStr,
		Lng:     lngStr,
		Address: address,
		Image:   header.Size,
	}

	if err := h.Validate(op, createReportDTO); err != nil {
		api.SendError(w, err)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 32)
	if err != nil {
		exc := api.NewSingleMessageException(
			api.EINVALID,
			op,
			fmt.Sprintf("Invalid %s as argument for Latitude. Latitude must be float.", latStr),
			err,
		)
		api.SendError(w, exc)
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 32)
	if err != nil {
		exc := api.NewSingleMessageException(
			api.EINVALID,
			op,
			fmt.Sprintf("Invalid %s as argument for Longitude. Longitude must be float.", lngStr),
			err,
		)
		api.SendError(w, exc)
		return
	}

	report := &entity.Report{
		UserID:  userPayload.ID,
		Address: address,
		Note:    note,
		Location: &entity.Location{
			Lat: lat,
			Lng: lng,
		},
	}

	report, err = h.ReportService.Create(r.Context(), report, image, header)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusCreated, "Created", report).SendJSON(w)
}

func (h *ReportHandler) GetAllReport(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	lastseenIDStr := r.URL.Query().Get("lastseenid")

	var limit uint64
	var lastseenID uint64
	var err error
	const op = "ReportHandler.GetAllReport"
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			exc := api.NewSingleMessageException(
				api.EINVALID,
				op,
				"Invalid limit argument",
				err,
			)
			api.SendError(w, exc)
			return
		}
	}

	if lastseenIDStr != "" {
		lastseenID, err = strconv.ParseUint(lastseenIDStr, 10, 64)
		if err != nil {
			exc := api.NewSingleMessageException(
				api.EINVALID,
				op,
				"Invalid lastseenid argument",
				err,
			)
			api.SendError(w, exc)
			return
		}
	}

	pagination := &model.Pagination{
		Limit:      limit,
		LastseenID: lastseenID,
	}
	reports, err := h.ReportService.GetAll(r.Context(), pagination)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusOK, "OK", reports).SendJSON(w)
}

func (h *ReportHandler) GetAllUserReport(w http.ResponseWriter, r *http.Request) {
	const op = "ReportHandler.GetAllUserReport"
	userPayload, err := api.UserPayloadFromContext(op, r)
	if err != nil {
		api.SendError(w, err)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	lastseenIDStr := r.URL.Query().Get("lastseenid")

	var limit uint64
	var lastseenID uint64
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			exc := api.NewSingleMessageException(
				api.EINVALID,
				op,
				"invalid limit argument",
				err,
			)
			api.SendError(w, exc)
			return
		}
	}

	if lastseenIDStr != "" {
		lastseenID, err = strconv.ParseUint(lastseenIDStr, 10, 64)
		if err != nil {
			exc := api.NewSingleMessageException(
				api.EINVALID,
				op,
				"invalid lastseenid argument",
				err,
			)
			api.SendError(w, exc)
			return
		}
	}

	pagination := &model.Pagination{
		Limit:      limit,
		LastseenID: lastseenID,
	}
	reports, err := h.ReportService.GetAllByUserID(r.Context(), userPayload.ID, pagination)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusOK, "OK", reports).SendJSON(w)
}

func (h *ReportHandler) UpdateReport(w http.ResponseWriter, r *http.Request) {
	const op = "ReportHandler.UpdateReport"
	userPayload, err := api.UserPayloadFromContext(op, r)
	if err != nil {
		api.SendError(w, err)
		return
	}

	if userPayload.Role != "ADMIN" {
		exc := api.NewSingleMessageException(
			api.EFORBIDDEN,
			op,
			"Forbidden",
			errors.New("trying to access admin endpoint without valid credential"),
		)
		api.SendError(w, exc)
		return
	}

	reportIDParam := chi.URLParam(r, "reportID")
	reportID, err := strconv.Atoi(reportIDParam)
	if err != nil {
		exc := api.NewSingleMessageException(
			api.EINVALID,
			op,
			"Invalid report id",
			err,
		)
		api.SendError(w, exc)
		return
	}

	updateReportDTO := new(model.UpdateReportDTO)
	if err := api.Bind(r.Body, updateReportDTO); err != nil {
		api.SendError(w, err)
		return
	}

	if err := h.Validate(op, updateReportDTO); err != nil {
		api.SendError(w, err)
		return
	}

	report, err := h.ReportService.Update(r.Context(), updateReportDTO.Status, reportID)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.NewResponse(http.StatusOK, "OK", report).SendJSON(w)
}
