package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

type TimeSlotModel struct {
	DB *sql.DB
}

type TimeSlot struct {
	ID            int64     `json:"timeslot_id"`
	ShiftID       int64     `json:"shift_id"`
	StartDateTime time.Time `json:"start_datetime"`
	EndDateTime   time.Time `json:"end_datetime"`
	IsPeak        bool      `json:"is_peak"`
}

func (m TimeSlotModel) Insert(timeSlot *TimeSlot) error {
	query := `
		INSERT INTO time_slots (shift_id, start_datetime, end_datetime, is_peak)
		VALUES ($1, $2, $3, $4)
		RETURNING timeslot_id
	`

	args := []any{
		timeSlot.ShiftID,
		timeSlot.StartDateTime,
		timeSlot.EndDateTime,
		timeSlot.IsPeak,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&timeSlot.ID)
}

func (m TimeSlotModel) Get(id int64) (*TimeSlot, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT timeslot_id, shift_id, start_datetime, end_datetime, is_peak
		FROM time_slots
		WHERE timeslot_id = $1
	`

	var timeSlot TimeSlot

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&timeSlot.ID,
		&timeSlot.ShiftID,
		&timeSlot.StartDateTime,
		&timeSlot.EndDateTime,
		&timeSlot.IsPeak,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &timeSlot, nil
}

func (m TimeSlotModel) Update(timeSlot *TimeSlot) error {
	query := `
		UPDATE time_slots
		SET shift_id = $1,
		    start_datetime = $2,
		    end_datetime = $3,
		    is_peak = $4
		WHERE timeslot_id = $5
		RETURNING timeslot_id
	`

	args := []any{
		timeSlot.ShiftID,
		timeSlot.StartDateTime,
		timeSlot.EndDateTime,
		timeSlot.IsPeak,
		timeSlot.ID,
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

func (m TimeSlotModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM time_slots
		WHERE timeslot_id = $1
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

func (m TimeSlotModel) GetAll(filters Filters) ([]*TimeSlot, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), timeslot_id, shift_id, start_datetime, end_datetime, is_peak
		FROM time_slots
		ORDER BY %s %s, timeslot_id ASC
        LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	timeSlots := []*TimeSlot{}

	for rows.Next() {
		var ts TimeSlot
		err := rows.Scan(
			&totalRecords,
			&ts.ID,
			&ts.ShiftID,
			&ts.StartDateTime,
			&ts.EndDateTime,
			&ts.IsPeak,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		timeSlots = append(timeSlots, &ts)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return timeSlots, metadata, nil
}

func ValidateTimeSlot(v *validator.Validator, ts *TimeSlot) {
	v.Check(ts.ShiftID > 0, "shift_id", "must be provided")
	v.Check(!ts.StartDateTime.IsZero(), "start_datetime", "must be provided")
	v.Check(!ts.EndDateTime.IsZero(), "end_datetime", "must be provided")
	v.Check(ts.EndDateTime.After(ts.StartDateTime), "end_datetime", "must be after start_datetime")
}
