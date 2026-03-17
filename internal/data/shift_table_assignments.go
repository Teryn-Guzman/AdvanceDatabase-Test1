package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ShiftTableAssignmentModel struct {
	DB *sql.DB
}

type ShiftTableAssignment struct {
	ShiftID int64 `json:"shift_id"`
	TableID int64 `json:"table_id"`
	StaffID int64 `json:"staff_id"`
}

func (m ShiftTableAssignmentModel) Insert(assignment *ShiftTableAssignment) error {
	query := `
		INSERT INTO shift_table_assignments (shift_id, table_id, staff_id)
		VALUES ($1, $2, $3)
	`

	args := []any{
		assignment.ShiftID,
		assignment.TableID,
		assignment.StaffID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m ShiftTableAssignmentModel) Get(shiftID, tableID int64) (*ShiftTableAssignment, error) {
	if shiftID < 1 || tableID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT shift_id, table_id, staff_id
		FROM shift_table_assignments
		WHERE shift_id = $1 AND table_id = $2
	`

	var assignment ShiftTableAssignment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, shiftID, tableID).Scan(
		&assignment.ShiftID,
		&assignment.TableID,
		&assignment.StaffID,
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

func (m ShiftTableAssignmentModel) Update(assignment *ShiftTableAssignment) error {
	query := `
		UPDATE shift_table_assignments
		SET staff_id = $1
		WHERE shift_id = $2 AND table_id = $3
	`

	args := []any{
		assignment.StaffID,
		assignment.ShiftID,
		assignment.TableID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
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

func (m ShiftTableAssignmentModel) Delete(shiftID, tableID int64) error {
	if shiftID < 1 || tableID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM shift_table_assignments
		WHERE shift_id = $1 AND table_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, shiftID, tableID)
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

func (m ShiftTableAssignmentModel) GetByShift(shiftID int64) ([]*ShiftTableAssignment, error) {
	if shiftID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT shift_id, table_id, staff_id
		FROM shift_table_assignments
		WHERE shift_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, shiftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*ShiftTableAssignment{}

	for rows.Next() {
		var a ShiftTableAssignment
		err := rows.Scan(&a.ShiftID, &a.TableID, &a.StaffID)
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

func (m ShiftTableAssignmentModel) GetAll(filters Filters) ([]*ShiftTableAssignment, Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), shift_id, table_id, staff_id
		FROM shift_table_assignments
		ORDER BY shift_id ASC, table_id ASC
        LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	assignments := []*ShiftTableAssignment{}

	for rows.Next() {
		var a ShiftTableAssignment
		err := rows.Scan(
			&totalRecords,
			&a.ShiftID,
			&a.TableID,
			&a.StaffID,
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
