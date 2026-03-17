package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

type ReservationModel struct {
	DB *sql.DB
}

type Reservation struct {
	ID          int64     `json:"reservation_id"`
	CustomerID  int64     `json:"customer_id"`
	TimeslotID  int64     `json:"timeslot_id"`
	PartySize   int       `json:"party_size"`
	Status      string    `json:"status"`
	IsWalkIn    bool      `json:"is_walk_in"`
	CreatedAt   time.Time `json:"created_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
	Notes       string    `json:"notes"`
}

func (m ReservationModel) Insert(reservation *Reservation) error {
	query := `
		INSERT INTO reservations (customer_id, timeslot_id, party_size, status, is_walk_in, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING reservation_id, created_at
	`

	args := []any{
		reservation.CustomerID,
		reservation.TimeslotID,
		reservation.PartySize,
		reservation.Status,
		reservation.IsWalkIn,
		reservation.Notes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&reservation.ID, &reservation.CreatedAt)
}

func (m ReservationModel) Get(id int64) (*Reservation, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT reservation_id, customer_id, timeslot_id, party_size, status, is_walk_in, created_at, cancelled_at, notes
		FROM reservations
		WHERE reservation_id = $1
	`

	var reservation Reservation

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.CustomerID,
		&reservation.TimeslotID,
		&reservation.PartySize,
		&reservation.Status,
		&reservation.IsWalkIn,
		&reservation.CreatedAt,
		&reservation.CancelledAt,
		&reservation.Notes,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &reservation, nil
}

func (m ReservationModel) Update(reservation *Reservation) error {
	query := `
		UPDATE reservations
		SET customer_id = $1,
		    timeslot_id = $2,
		    party_size = $3,
		    status = $4,
		    is_walk_in = $5,
		    cancelled_at = $6,
		    notes = $7
		WHERE reservation_id = $8
		RETURNING reservation_id
	`

	args := []any{
		reservation.CustomerID,
		reservation.TimeslotID,
		reservation.PartySize,
		reservation.Status,
		reservation.IsWalkIn,
		reservation.CancelledAt,
		reservation.Notes,
		reservation.ID,
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

func (m ReservationModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM reservations
		WHERE reservation_id = $1
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

func (m ReservationModel) GetAll(filters Filters) ([]*Reservation, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), reservation_id, customer_id, timeslot_id, party_size, status, is_walk_in, created_at, cancelled_at, notes
		FROM reservations
		ORDER BY %s %s, reservation_id ASC
        LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	reservations := []*Reservation{}

	for rows.Next() {
		var r Reservation
		err := rows.Scan(
			&totalRecords,
			&r.ID,
			&r.CustomerID,
			&r.TimeslotID,
			&r.PartySize,
			&r.Status,
			&r.IsWalkIn,
			&r.CreatedAt,
			&r.CancelledAt,
			&r.Notes,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reservations = append(reservations, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return reservations, metadata, nil
}

func ValidateReservation(v *validator.Validator, r *Reservation) {
	v.Check(r.CustomerID > 0, "customer_id", "must be provided")
	v.Check(r.TimeslotID > 0, "timeslot_id", "must be provided")
	v.Check(r.PartySize > 0, "party_size", "must be greater than 0")
	v.Check(r.Status != "", "status", "must be provided")
	v.Check(r.Status == "confirmed" || r.Status == "cancelled" || r.Status == "no_show" || r.Status == "completed", 
		"status", "must be one of: confirmed, cancelled, no_show, completed")
}
