package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

type ShiftModel struct {
	DB *sql.DB
}

type Shift struct {
	ID        int64     `json:"shift_id"`
	ShiftName string    `json:"shift_name"`
	ShiftDate time.Time `json:"shift_date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (m ShiftModel) Insert(shift *Shift) error {
	query := `
		INSERT INTO shifts (shift_name, shift_date, start_time, end_time)
		VALUES ($1, $2, $3, $4)
		RETURNING shift_id
	`

	args := []any{
		shift.ShiftName,
		shift.ShiftDate,
		shift.StartTime,
		shift.EndTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&shift.ID)
}

func (m ShiftModel) Get(id int64) (*Shift, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT shift_id, shift_name, shift_date, start_time, end_time
		FROM shifts
		WHERE shift_id = $1
	`

	var shift Shift

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&shift.ID,
		&shift.ShiftName,
		&shift.ShiftDate,
		&shift.StartTime,
		&shift.EndTime,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &shift, nil
}

func (m ShiftModel) Update(shift *Shift) error {
	query := `
		UPDATE shifts
		SET shift_name = $1,
		    shift_date = $2,
		    start_time = $3,
		    end_time = $4
		WHERE shift_id = $5
		RETURNING shift_id
	`

	args := []any{
		shift.ShiftName,
		shift.ShiftDate,
		shift.StartTime,
		shift.EndTime,
		shift.ID,
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

func (m ShiftModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM shifts
		WHERE shift_id = $1
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

func (m ShiftModel) GetAll(filters Filters) ([]*Shift, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), shift_id, shift_name, shift_date, start_time, end_time
		FROM shifts
		ORDER BY %s %s, shift_id ASC
        LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	shifts := []*Shift{}

	for rows.Next() {
		var s Shift
		err := rows.Scan(
			&totalRecords,
			&s.ID,
			&s.ShiftName,
			&s.ShiftDate,
			&s.StartTime,
			&s.EndTime,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		shifts = append(shifts, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return shifts, metadata, nil
}

func ValidateShift(v *validator.Validator, s *Shift) {
	if s.ShiftName != "" {
		v.Check(len(s.ShiftName) <= 100, "shift_name", "must not exceed 100 characters")
	}

	v.Check(!s.ShiftDate.IsZero(), "shift_date", "must be provided")
	v.Check(!s.StartTime.IsZero(), "start_time", "must be provided")
	v.Check(!s.EndTime.IsZero(), "end_time", "must be provided")
	v.Check(s.EndTime.After(s.StartTime), "end_time", "must be after start_time")
}
