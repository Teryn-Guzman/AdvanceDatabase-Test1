package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

func (a *applicationDependencies) createShiftHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		ShiftName string `json:"shift_name"`
		ShiftDate string `json:"shift_date"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	shiftDate, _ := time.Parse("2006-01-02", input.ShiftDate)
	startTime, _ := time.Parse("15:04:05", input.StartTime)
	endTime, _ := time.Parse("15:04:05", input.EndTime)

	shift := data.Shift{
		ShiftName: input.ShiftName,
		ShiftDate: shiftDate,
		StartTime: startTime,
		EndTime:   endTime,
	}

	v := validator.New()
	data.ValidateShift(v, &shift)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.shiftModel.Insert(&shift)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/shifts/%d", shift.ID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"shift": shift}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayShiftHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	shift, err := a.shiftModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"shift": shift,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateShiftHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	shift, err := a.shiftModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		ShiftName *string `json:"shift_name"`
		ShiftDate *string `json:"shift_date"`
		StartTime *string `json:"start_time"`
		EndTime   *string `json:"end_time"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.ShiftName != nil {
		shift.ShiftName = *input.ShiftName
	}
	if input.ShiftDate != nil {
		shiftDate, _ := time.Parse("2006-01-02", *input.ShiftDate)
		shift.ShiftDate = shiftDate
	}
	if input.StartTime != nil {
		startTime, _ := time.Parse("15:04:05", *input.StartTime)
		shift.StartTime = startTime
	}
	if input.EndTime != nil {
		endTime, _ := time.Parse("15:04:05", *input.EndTime)
		shift.EndTime = endTime
	}

	v := validator.New()
	data.ValidateShift(v, shift)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.shiftModel.Update(shift)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	dataEnvelope := envelope{
		"shift": shift,
	}

	err = a.writeJSON(w, http.StatusOK, dataEnvelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteShiftHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.shiftModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"message": "shift successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listShiftsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var queryParametersData struct {
		data.Filters
	}

	queryParameters := r.URL.Query()

	v := validator.New()

	queryParametersData.Filters.Page = a.getSingleIntegerParameter(
		queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(
		queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(
		queryParameters, "sort", "shift_id")

	queryParametersData.Filters.SortSafeList = []string{"shift_id", "shift_date",
		"-shift_id", "-shift_date"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	shifts, metadata, err := a.shiftModel.GetAll(queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"shifts":    shifts,
		"@metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
