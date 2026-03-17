package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

// Handlers for managing reservations
func (a *applicationDependencies) createReservationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		CustomerID  int64  `json:"customer_id"`
		TimeslotID  int64  `json:"timeslot_id"`
		PartySize   int    `json:"party_size"`
		Status      string `json:"status"`
		IsWalkIn    bool   `json:"is_walk_in"`
		Notes       string `json:"notes"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	reservation := data.Reservation{
		CustomerID: input.CustomerID,
		TimeslotID: input.TimeslotID,
		PartySize:  input.PartySize,
		Status:     input.Status,
		IsWalkIn:   input.IsWalkIn,
		Notes:      input.Notes,
	}

	v := validator.New()
	data.ValidateReservation(v, &reservation)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.reservationModel.Insert(&reservation)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/reservations/%d", reservation.ID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"reservation": reservation}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Additional handlers for displaying, updating, deleting, and listing reservations would follow a similar pattern to the create handler, with appropriate adjustments for the specific operation being performed.
func (a *applicationDependencies) displayReservationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	reservation, err := a.reservationModel.Get(id)
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
		"reservation": reservation,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Similar handlers for updating, deleting, and listing reservations would be implemented here, following the same structure and error handling patterns as the display handler.
func (a *applicationDependencies) updateReservationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	reservation, err := a.reservationModel.Get(id)
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
		CustomerID  *int64  `json:"customer_id"`
		TimeslotID  *int64  `json:"timeslot_id"`
		PartySize   *int    `json:"party_size"`
		Status      *string `json:"status"`
		IsWalkIn    *bool   `json:"is_walk_in"`
		Notes       *string `json:"notes"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.CustomerID != nil {
		reservation.CustomerID = *input.CustomerID
	}
	if input.TimeslotID != nil {
		reservation.TimeslotID = *input.TimeslotID
	}
	if input.PartySize != nil {
		reservation.PartySize = *input.PartySize
	}
	if input.Status != nil {
		reservation.Status = *input.Status
	}
	if input.IsWalkIn != nil {
		reservation.IsWalkIn = *input.IsWalkIn
	}
	if input.Notes != nil {
		reservation.Notes = *input.Notes
	}

	v := validator.New()
	data.ValidateReservation(v, reservation)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.reservationModel.Update(reservation)
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
		"reservation": reservation,
	}

	err = a.writeJSON(w, http.StatusOK, dataEnvelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// Similar handlers for deleting and listing reservations would be implemented here, following the same structure and error handling patterns as the display handler.
func (a *applicationDependencies) deleteReservationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.reservationModel.Delete(id)
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
		"message": "reservation successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Handler for listing reservations with pagination and sorting
func (a *applicationDependencies) listReservationsHandler(
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
		queryParameters, "sort", "reservation_id")

	queryParametersData.Filters.SortSafeList = []string{"reservation_id", "customer_id", "status",
		"-reservation_id", "-customer_id", "-status"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	reservations, metadata, err := a.reservationModel.GetAll(queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"reservations": reservations,
		"@metadata":    metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
