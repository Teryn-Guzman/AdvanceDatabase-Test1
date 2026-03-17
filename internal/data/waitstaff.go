package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

type WaitstaffModel struct {
	DB *sql.DB
}

type Waitstaff struct {
	ID        int64     `json:"staff_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	HireDate  time.Time `json:"hire_date"`
	IsActive  bool      `json:"is_active"`
}

func (m WaitstaffModel) Insert(staff *Waitstaff) error {
	query := `
		INSERT INTO waitstaff (first_name, last_name, hire_date, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING staff_id
	`

	args := []any{
		staff.FirstName,
		staff.LastName,
		staff.HireDate,
		staff.IsActive,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&staff.ID)
}

func (m WaitstaffModel) Get(id int64) (*Waitstaff, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT staff_id, first_name, last_name, hire_date, is_active
		FROM waitstaff
		WHERE staff_id = $1
	`

	var staff Waitstaff

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&staff.ID,
		&staff.FirstName,
		&staff.LastName,
		&staff.HireDate,
		&staff.IsActive,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &staff, nil
}

func (m WaitstaffModel) Update(staff *Waitstaff) error {
	query := `
		UPDATE waitstaff
		SET first_name = $1,
		    last_name = $2,
		    hire_date = $3,
		    is_active = $4
		WHERE staff_id = $5
		RETURNING staff_id
	`

	args := []any{
		staff.FirstName,
		staff.LastName,
		staff.HireDate,
		staff.IsActive,
		staff.ID,
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

func (m WaitstaffModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM waitstaff
		WHERE staff_id = $1
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

func (m WaitstaffModel) GetAll(filters Filters) ([]*Waitstaff, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), staff_id, first_name, last_name, hire_date, is_active
		FROM waitstaff
		ORDER BY %s %s, staff_id ASC
        LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	staff := []*Waitstaff{}

	for rows.Next() {
		var w Waitstaff
		err := rows.Scan(
			&totalRecords,
			&w.ID,
			&w.FirstName,
			&w.LastName,
			&w.HireDate,
			&w.IsActive,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		staff = append(staff, &w)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return staff, metadata, nil
}

func ValidateWaitstaff(v *validator.Validator, w *Waitstaff) {
	v.Check(w.FirstName != "", "first_name", "must be provided")
	v.Check(len(w.FirstName) <= 100, "first_name", "must not exceed 100 characters")

	v.Check(w.LastName != "", "last_name", "must be provided")
	v.Check(len(w.LastName) <= 100, "last_name", "must not exceed 100 characters")

	v.Check(!w.HireDate.IsZero(), "hire_date", "must be provided")
}
