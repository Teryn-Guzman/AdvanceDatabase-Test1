package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

func (a *applicationDependencies) createShiftTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		ShiftID int64 `json:"shift_id"`
		TableID int64 `json:"table_id"`
		StaffID int64 `json:"staff_id"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	assignment := data.ShiftTableAssignment{
		ShiftID: input.ShiftID,
		TableID: input.TableID,
		StaffID: input.StaffID,
	}

	err = a.shiftTableAssignmentModel.Insert(&assignment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/shift-table-assignments/%d/%d", assignment.ShiftID, assignment.TableID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"assignment": assignment}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayShiftTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	shiftID, err := strconv.ParseInt(r.URL.Query().Get("shift_id"), 10, 64)
	if err != nil || shiftID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid shift_id"))
		return
	}

	tableID, err := strconv.ParseInt(r.URL.Query().Get("table_id"), 10, 64)
	if err != nil || tableID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid table_id"))
		return
	}

	assignment, err := a.shiftTableAssignmentModel.Get(shiftID, tableID)
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

func (a *applicationDependencies) updateShiftTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	shiftID, err := strconv.ParseInt(r.URL.Query().Get("shift_id"), 10, 64)
	if err != nil || shiftID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid shift_id"))
		return
	}

	tableID, err := strconv.ParseInt(r.URL.Query().Get("table_id"), 10, 64)
	if err != nil || tableID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid table_id"))
		return
	}

	assignment, err := a.shiftTableAssignmentModel.Get(shiftID, tableID)
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
		StaffID *int64 `json:"staff_id"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.StaffID != nil {
		assignment.StaffID = *input.StaffID
	}

	err = a.shiftTableAssignmentModel.Update(assignment)
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
		"assignment": assignment,
	}

	err = a.writeJSON(w, http.StatusOK, dataEnvelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteShiftTableAssignmentHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	shiftID, err := strconv.ParseInt(r.URL.Query().Get("shift_id"), 10, 64)
	if err != nil || shiftID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid shift_id"))
		return
	}

	tableID, err := strconv.ParseInt(r.URL.Query().Get("table_id"), 10, 64)
	if err != nil || tableID < 1 {
		a.badRequestResponse(w, r, fmt.Errorf("invalid table_id"))
		return
	}

	err = a.shiftTableAssignmentModel.Delete(shiftID, tableID)
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
		"message": "shift table assignment successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listShiftTableAssignmentsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var queryParametersData struct {
		data.Filters
	}

	queryParameters := r.URL.Query()

	v := validator.New()
	if filters, ok := queryParameters["filters"]; ok && len(filters) > 0 {
		// parse filters if needed
	}

	queryParametersData.Filters.Page = a.getSingleIntegerParameter(
		queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(
		queryParameters, "page_size", 10, v)

	queryParametersData.Filters.Sort = a.getSingleQueryParameter(
		queryParameters, "sort", "shift_id")

	queryParametersData.Filters.SortSafeList = []string{"shift_id", "table_id",
		"-shift_id", "-table_id"}

	assignments, metadata, err := a.shiftTableAssignmentModel.GetAll(queryParametersData.Filters)
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
