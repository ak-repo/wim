package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/model"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrCustomerNotFound = errors.New("customer not found")

type CustomerRepository interface {
	Create(ctx context.Context, customer *model.CustomerRequest) (int, error)
	GetByID(ctx context.Context, customerID int) (*model.CustomerDTO, error)
	GetByEmail(ctx context.Context, email string) (*model.CustomerDTO, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, customerID int, customer *model.CustomerRequest) error
	Delete(ctx context.Context, customerID int) error
	List(ctx context.Context, params *model.CustomerParams) (model.CustomerDTOs, error)
	Count(ctx context.Context, params *model.CustomerParams) (int, error)
}

type customerRepository struct {
	db *db.DB
}

func NewCustomerRepository(database *db.DB) CustomerRepository {
	return &customerRepository{db: database}
}

func (r *customerRepository) Create(ctx context.Context, customer *model.CustomerRequest) (int, error) {
	query := `
		INSERT INTO customers (
			ref_code, name, email, contact, address, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.Pool.QueryRow(ctx, query,
		customer.RefCode,
		customer.Name,
		customer.Email,
		customer.Contact,
		customer.Address,
		customer.IsActive,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrAlreadyExists
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create customer")
	}

	return id, nil
}

func (r *customerRepository) GetByID(ctx context.Context, customerID int) (*model.CustomerDTO, error) {
	return scanCustomer(ctx, r.db, `
		SELECT id, ref_code, name, email, contact, address, is_active, created_at, updated_at, deleted_at
		FROM customers WHERE id = $1
	`, customerID)
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*model.CustomerDTO, error) {
	return scanCustomer(ctx, r.db, `
		SELECT id, ref_code, name, email, contact, address, is_active, created_at, updated_at, deleted_at
		FROM customers WHERE email = $1
	`, email)
}

func (r *customerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check customer by email")
	}

	return exists, nil
}

func (r *customerRepository) Update(ctx context.Context, customerID int, customer *model.CustomerRequest) error {
	query := `
		UPDATE customers
		SET name = COALESCE($2, name),
			email = COALESCE($3, email),
			contact = COALESCE($4, contact),
			address = COALESCE($5, address),
			is_active = COALESCE($6, is_active),
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		customerID,
		customer.Name,
		customer.Email,
		customer.Contact,
		customer.Address,
		customer.IsActive,
	)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update customer")
	}

	if result.RowsAffected() == 0 {
		return ErrCustomerNotFound
	}

	return nil
}

func (r *customerRepository) Delete(ctx context.Context, customerID int) error {
	result, err := r.db.Pool.Exec(ctx, `DELETE FROM customers WHERE id = $1`, customerID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete customer")
	}

	if result.RowsAffected() == 0 {
		return ErrCustomerNotFound
	}

	return nil
}

func (r *customerRepository) List(ctx context.Context, params *model.CustomerParams) (model.CustomerDTOs, error) {
	var (
		args       []any
		conditions []string
	)

	query := `
		SELECT id, ref_code, name, email, contact, address, is_active, created_at, updated_at, deleted_at
		FROM customers
	`

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, offset)

	return scanCustomers(ctx, r.db, query, args...)
}

func (r *customerRepository) Count(ctx context.Context, params *model.CustomerParams) (int, error) {
	var (
		count      int
		args       []any
		conditions []string
	)

	query := `SELECT COUNT(*) FROM customers`

	if params.Active != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.Active)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count customers")
	}

	return count, nil
}

func scanCustomer(ctx context.Context, database *db.DB, query string, args ...any) (*model.CustomerDTO, error) {
	var row model.CustomerDTO
	var isActive sql.NullBool
	var createdAt, updatedAt, deletedAt sql.NullTime

	err := database.Pool.QueryRow(ctx, query, args...).Scan(
		&row.ID,
		&row.RefCode,
		&row.Name,
		&row.Email,
		&row.Contact,
		&row.Address,
		&isActive,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load customer")
	}
	row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

	return &row, nil
}

func scanCustomers(ctx context.Context, database *db.DB, query string, args ...any) (model.CustomerDTOs, error) {
	rows, err := database.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list customers")
	}
	defer rows.Close()

	var customers model.CustomerDTOs

	for rows.Next() {
		var row model.CustomerDTO
		var isActive sql.NullBool
		var createdAt, updatedAt, deletedAt sql.NullTime

		err := rows.Scan(
			&row.ID,
			&row.RefCode,
			&row.Name,
			&row.Email,
			&row.Contact,
			&row.Address,
			&isActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to scan customer row")
		}
		row.ApplyNullScalars(isActive, createdAt, updatedAt, deletedAt)

		customers = append(customers, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to iterate customers")
	}

	return customers, nil
}
