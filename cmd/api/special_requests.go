package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

func (a *applicationDependencies) createSpecialRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		ReservationID int64  `json:"reservation_id"`
		RequestType   string `json:"request_type"`
		Description   string `json:"description"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	request := data.SpecialRequest{
		ReservationID: input.ReservationID,
		RequestType:   input.RequestType,
		Description:   input.Description,
	}

	v := validator.New()
	data.ValidateSpecialRequest(v, &request)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.specialRequestModel.Insert(&request)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/special-requests/%d", request.ID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"request": request}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displaySpecialRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	request, err := a.specialRequestModel.Get(id)
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
		"request": request,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateSpecialRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	request, err := a.specialRequestModel.Get(id)
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
		ReservationID *int64  `json:"reservation_id"`
		RequestType   *string `json:"request_type"`
		Description   *string `json:"description"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.ReservationID != nil {
		request.ReservationID = *input.ReservationID
	}
	if input.RequestType != nil {
		request.RequestType = *input.RequestType
	}
	if input.Description != nil {
		request.Description = *input.Description
	}

	v := validator.New()
	data.ValidateSpecialRequest(v, request)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.specialRequestModel.Update(request)
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
		"request": request,
	}

	err = a.writeJSON(w, http.StatusOK, dataEnvelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteSpecialRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.specialRequestModel.Delete(id)
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
		"message": "special request successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listSpecialRequestsHandler(
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
		queryParameters, "sort", "request_id")

	queryParametersData.Filters.SortSafeList = []string{"request_id", "reservation_id",
		"-request_id", "-reservation_id"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	requests, metadata, err := a.specialRequestModel.GetAll(queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"requests":  requests,
		"@metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
