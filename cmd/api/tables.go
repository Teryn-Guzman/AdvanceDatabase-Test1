package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

func (a *applicationDependencies) createTableHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input struct {
		TableNumber string `json:"table_number"`
		Capacity    int    `json:"capacity"`
		Location    string `json:"location"`
		IsActive    bool   `json:"is_active"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	table := data.Table{
		TableNumber: input.TableNumber,
		Capacity:    input.Capacity,
		Location:    input.Location,
		IsActive:    input.IsActive,
	}

	v := validator.New()
	data.ValidateTable(v, &table)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.tableModel.Insert(&table)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tables/%d", table.ID))

	err = a.writeJSON(w, http.StatusCreated,
		envelope{"table": table}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayTableHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	table, err := a.tableModel.Get(id)
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
		"table": table,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateTableHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	table, err := a.tableModel.Get(id)
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
		TableNumber *string `json:"table_number"`
		Capacity    *int    `json:"capacity"`
		Location    *string `json:"location"`
		IsActive    *bool   `json:"is_active"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.TableNumber != nil {
		table.TableNumber = *input.TableNumber
	}
	if input.Capacity != nil {
		table.Capacity = *input.Capacity
	}
	if input.Location != nil {
		table.Location = *input.Location
	}
	if input.IsActive != nil {
		table.IsActive = *input.IsActive
	}

	v := validator.New()
	data.ValidateTable(v, table)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.tableModel.Update(table)
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
		"table": table,
	}

	err = a.writeJSON(w, http.StatusOK, dataEnvelope, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteTableHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.tableModel.Delete(id)
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
		"message": "table successfully deleted",
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listTablesHandler(
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
		queryParameters, "sort", "table_id")

	queryParametersData.Filters.SortSafeList = []string{"table_id", "table_number",
		"-table_id", "-table_number"}

	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	tables, metadata, err := a.tableModel.GetAll(queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"tables":    tables,
		"@metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
