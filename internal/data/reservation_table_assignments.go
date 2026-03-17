package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ReservationTableAssignmentModel struct {
	DB *sql.DB
}

type ReservationTableAssignment struct {
	ReservationID int64 `json:"reservation_id"`
	TableID       int64 `json:"table_id"`
}

func (m ReservationTableAssignmentModel) Insert(assignment *ReservationTableAssignment) error {
	query := `
		INSERT INTO reservation_table_assignments (reservation_id, table_id)
		VALUES ($1, $2)
	`

	args := []any{
		assignment.ReservationID,
		assignment.TableID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m ReservationTableAssignmentModel) Get(reservationID, tableID int64) (*ReservationTableAssignment, error) {
	if reservationID < 1 || tableID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT reservation_id, table_id
		FROM reservation_table_assignments
		WHERE reservation_id = $1 AND table_id = $2
	`

	var assignment ReservationTableAssignment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, reservationID, tableID).Scan(
		&assignment.ReservationID,
		&assignment.TableID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &assignment, nil
}

func (m ReservationTableAssignmentModel) Delete(reservationID, tableID int64) error {
	if reservationID < 1 || tableID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM reservation_table_assignments
		WHERE reservation_id = $1 AND table_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, reservationID, tableID)
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

func (m ReservationTableAssignmentModel) GetByReservation(reservationID int64) ([]*ReservationTableAssignment, error) {
	if reservationID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT reservation_id, table_id
		FROM reservation_table_assignments
		WHERE reservation_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*ReservationTableAssignment{}

	for rows.Next() {
		var a ReservationTableAssignment
		err := rows.Scan(&a.ReservationID, &a.TableID)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (m ReservationTableAssignmentModel) GetAll(filters Filters) ([]*ReservationTableAssignment, Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), reservation_id, table_id
		FROM reservation_table_assignments
		ORDER BY reservation_id ASC, table_id ASC
        LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	assignments := []*ReservationTableAssignment{}

	for rows.Next() {
		var a ReservationTableAssignment
		err := rows.Scan(
			&totalRecords,
			&a.ReservationID,
			&a.TableID,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		assignments = append(assignments, &a)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return assignments, metadata, nil
}
