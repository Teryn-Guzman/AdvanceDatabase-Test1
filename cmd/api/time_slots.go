package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

func (a *applicationDependencies) createTimeSlotHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		ShiftID       int64  `json:"shift_id"`
		StartDateTime string `json:"start_datetime"`
		EndDateTime   string `json:"end_datetime"`
		IsPeak        bool   `json:"is_peak"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	startDateTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", input.StartDateTime)
	endDateTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", input.EndDateTime)

	timeSlot := data.TimeSlot{
		ShiftID:       input.ShiftID,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		IsPeak:        input.IsPeak,
	}

	v := validator.New()
	data.ValidateTimeSlot(v, &timeSlot)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.timeSlotModel.Insert(&timeSlot)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/time-slots/%d", timeSlot.ID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"time_slot": timeSlot}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayTimeSlotHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	timeSlot, err := a.timeSlotModel.Get(id)
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
		"time_slot": timeSlot,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateTimeSlotHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	timeSlot, err := a.timeSlotModel.Get(id)
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
		ShiftID       *int64  `json:"shift_id"`
		StartDateTime *string `json:"start_datetime"`
		EndDateTime   *string `json:"end_datetime"`
		IsPeak        *bool   `json:"is_peak"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.ShiftID != nil {
		timeSlot.ShiftID = *input.ShiftID
	}
	if input.StartDateTime != nil {
		startDateTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", *input.StartDateTime)
		timeSlot.StartDateTime = startDateTime
	}
	if input.EndDateTime != nil {
		endDateTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", *input.EndDateTime)
		timeSlot.EndDateTime = endDateTime
	}
	if input.IsPeak != nil {
		timeSlot.IsPeak = *input.IsPeak
	}

	v := validator.New()
	data.ValidateTimeSlot(v, timeSlot)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.timeSlotModel.Update(timeSlot)
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
		"time_slot": timeSlot,
	}

	err = a.writeJSON(w, http.StatusOK, dataEnvelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteTimeSlotHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.timeSlotModel.Delete(id)
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
		"message": "time slot successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listTimeSlotsHandler(
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
		queryParameters, "sort", "timeslot_id")

	queryParametersData.Filters.SortSafeList = []string{"timeslot_id", "shift_id",
		"-timeslot_id", "-shift_id"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	timeSlots, metadata, err := a.timeSlotModel.GetAll(queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"time_slots": timeSlots,
		"@metadata":  metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
