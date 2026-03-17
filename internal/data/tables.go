package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/validator"
)

type TableModel struct {
	DB *sql.DB
}

type Table struct {
	ID         int64  `json:"table_id"`
	TableNumber string `json:"table_number"`
	Capacity   int    `json:"capacity"`
	Location   string `json:"location"`
	IsActive   bool   `json:"is_active"`
}

func (m TableModel) Insert(table *Table) error {
	query := `
		INSERT INTO tables (table_number, capacity, location, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING table_id
	`

	args := []any{
		table.TableNumber,
		table.Capacity,
		table.Location,
		table.IsActive,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&table.ID)
}

func (m TableModel) Get(id int64) (*Table, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT table_id, table_number, capacity, location, is_active
		FROM tables
		WHERE table_id = $1
	`

	var table Table

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&table.ID,
		&table.TableNumber,
		&table.Capacity,
		&table.Location,
		&table.IsActive,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &table, nil
}

func (m TableModel) Update(table *Table) error {
	query := `
		UPDATE tables
		SET table_number = $1,
		    capacity = $2,
		    location = $3,
		    is_active = $4
		WHERE table_id = $5
		RETURNING table_id
	`

	args := []any{
		table.TableNumber,
		table.Capacity,
		table.Location,
		table.IsActive,
		table.ID,
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

func (m TableModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM tables
		WHERE table_id = $1
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

func (m TableModel) GetAll(filters Filters) ([]*Table, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), table_id, table_number, capacity, location, is_active
		FROM tables
		ORDER BY %s %s, table_id ASC
        LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	tables := []*Table{}

	for rows.Next() {
		var t Table
		err := rows.Scan(
			&totalRecords,
			&t.ID,
			&t.TableNumber,
			&t.Capacity,
			&t.Location,
			&t.IsActive,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		tables = append(tables, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return tables, metadata, nil
}

func ValidateTable(v *validator.Validator, t *Table) {
	v.Check(t.TableNumber != "", "table_number", "must be provided")
	v.Check(len(t.TableNumber) <= 20, "table_number", "must not exceed 20 characters")

	v.Check(t.Capacity > 0, "capacity", "must be greater than 0")
	v.Check(t.Capacity <= 100, "capacity", "must not exceed 100")

	if t.Location != "" {
		v.Check(len(t.Location) <= 100, "location", "must not exceed 100 characters")
	}
}
