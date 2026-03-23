package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

// Handlers for managing reservation table assignments
func (a *applicationDependencies) createReservationTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		ReservationID int64 `json:"reservation_id"`
		TableID       int64 `json:"table_id"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	assignment := data.ReservationTableAssignment{
		ReservationID: input.ReservationID,
		TableID:       input.TableID,
	}

	err = a.reservationTableAssignmentModel.Insert(&assignment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/reservation-table-assignments/%d/%d", assignment.ReservationID, assignment.TableID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"assignment": assignment}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Additional handlers for displaying, deleting, and listing reservation table assignments would follow a similar pattern to the shift table assignment handlers, with appropriate adjustments for the reservation context.
func (a *applicationDependencies) displayReservationTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	reservationID, err := strconv.ParseInt(r.URL.Query().Get("reservation_id"), 10, 64)
	if err != nil || reservationID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid reservation_id"))
		return
	}

	tableID, err := strconv.ParseInt(r.URL.Query().Get("table_id"), 10, 64)
	if err != nil || tableID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid table_id"))
		return
	}

	assignment, err := a.reservationTableAssignmentModel.Get(reservationID, tableID)
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
		"assignment": assignment,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Similar handlers for deleting and listing reservation table assignments would be implemented here, following the same structure and error handling patterns as the display handler.
func (a *applicationDependencies) deleteReservationTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	reservationID, err := strconv.ParseInt(r.URL.Query().Get("reservation_id"), 10, 64)
	if err != nil || reservationID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid reservation_id"))
		return
	}

	tableID, err := strconv.ParseInt(r.URL.Query().Get("table_id"), 10, 64)
	if err != nil || tableID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid table_id"))
		return
	}

	err = a.reservationTableAssignmentModel.Delete(reservationID, tableID)
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
		"message": "reservation table assignment successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// Handler for listing reservation table assignments with pagination and sorting
func (a *applicationDependencies) listReservationTableAssignmentsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	queryParameters := r.URL.Query()

	v := validator.New()
	if filters, ok := queryParameters["filters"]; ok && len(filters) > 0 {
		// parse filters if needed
	}

	var queryParametersData struct {
		data.Filters
	}

	queryParametersData.Filters.Page = a.getSingleIntegerParameter(
		queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(
		queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(
		queryParameters, "sort", "reservation_id")

	queryParametersData.Filters.SortSafeList = []string{"reservation_id", "table_id",
		"-reservation_id", "-table_id"}

	assignments, metadata, err := a.reservationTableAssignmentModel.GetAll(queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"assignments": assignments,
		"@metadata":   metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}