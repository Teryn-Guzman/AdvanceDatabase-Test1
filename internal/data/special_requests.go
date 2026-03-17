package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

type SpecialRequestModel struct {
	DB *sql.DB
}

type SpecialRequest struct {
	ID            int64  `json:"request_id"`
	ReservationID int64  `json:"reservation_id"`
	RequestType   string `json:"request_type"`
	Description   string `json:"description"`
}

func (m SpecialRequestModel) Insert(request *SpecialRequest) error {
	query := `
		INSERT INTO special_requests (reservation_id, request_type, description)
		VALUES ($1, $2, $3)
		RETURNING request_id
	`

	args := []any{
		request.ReservationID,
		request.RequestType,
		request.Description,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&request.ID)
}

func (m SpecialRequestModel) Get(id int64) (*SpecialRequest, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT request_id, reservation_id, request_type, description
		FROM special_requests
		WHERE request_id = $1
	`

	var request SpecialRequest

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.ReservationID,
		&request.RequestType,
		&request.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &request, nil
}

func (m SpecialRequestModel) Update(request *SpecialRequest) error {
	query := `
		UPDATE special_requests
		SET reservation_id = $1,
		    request_type = $2,
		    description = $3
		WHERE request_id = $4
		RETURNING request_id
	`

	args := []any{
		request.ReservationID,
		request.RequestType,
		request.Description,
		request.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int64
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m SpecialRequestModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM special_requests
		WHERE request_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m SpecialRequestModel) GetAll(filters Filters) ([]*SpecialRequest, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), request_id, reservation_id, request_type, description
		FROM special_requests
		ORDER BY %s %s, request_id ASC
        LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	requests := []*SpecialRequest{}

	for rows.Next() {
		var sr SpecialRequest
		err := rows.Scan(
			&totalRecords,
			&sr.ID,
			&sr.ReservationID,
			&sr.RequestType,
			&sr.Description,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		requests = append(requests, &sr)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return requests, metadata, nil
}

func ValidateSpecialRequest(v *validator.Validator, sr *SpecialRequest) {
	v.Check(sr.ReservationID > 0, "reservation_id", "must be provided")
	
	if sr.RequestType != "" {
		v.Check(len(sr.RequestType) <= 100, "request_type", "must not exceed 100 characters")
	}

	if sr.Description != "" {
		v.Check(len(sr.Description) <= 1000, "description", "must not exceed 1000 characters")
	}
}
