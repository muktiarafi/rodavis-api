package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/api"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/driver"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/entity"
)

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, e driver.Executor, user *entity.User) (*entity.User, error) {
	ctx, cancel := newDBContext(ctx)
	defer cancel()

	stmt := `INSERT INTO users (name, phone_number, email, password)
	VALUES ($1, $2, $3, $4)
	RETURNING *`

	const op = "UserRepositoryImpl.Create"
	newUser := new(entity.User)
	if err := e.QueryRowContext(
		ctx,
		stmt,
		user.Name,
		user.PhoneNumber,
		user.Email,
		user.Password,
	).Scan(
		&newUser.ID,
		&newUser.Name,
		&newUser.PhoneNumber,
		&newUser.Email,
		&newUser.Password,
		&newUser.Role,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_email_key" {
				return nil, api.NewSingleMessageException(
					api.ECONFLICT,
					op,
					"Email already taken",
					errors.New("trying to register with already taken email"),
				)
			}

			if pgerr.ConstraintName == "users_phone_number_key" {
				return nil, api.NewSingleMessageException(
					api.ECONFLICT,
					op,
					"Phonenumber already taken",
					errors.New("trying to register with already taken phonenumber"),
				)
			}
		}
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"r.SQL.QueryRowContext",
			err,
		)
	}

	return newUser, nil
}

func (r *UserRepositoryImpl) Get(ctx context.Context, e driver.Executor, userID int) (*entity.User, error) {
	ctx, cancel := newDBContext(ctx)
	defer cancel()

	stmt := `SELECT * FROM users
	WHERE id = $1`

	user := new(entity.User)
	if err := e.QueryRowContext(ctx, stmt, userID).Scan(
		&user.ID,
		&user.Name,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		const op = "UserRepositoryImpl.Get"
		if err == sql.ErrNoRows {
			return nil, api.NewSingleMessageException(
				api.ENOTFOUND,
				op,
				"User Not Found",
				err,
			)
		}
		return nil, api.NewExceptionWithSourceLocation(
			op,
			"r.SQL.QueryRowContext",
			err,
		)
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, e driver.Executor, email string) (*entity.User, error) {
	ctx, cancel := newDBContext(ctx)
	defer cancel()

	stmt := `SELECT * FROM users
	WHERE email = $1`

	user := new(entity.User)
	if err := e.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.Name,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		const op = "UserRepositoryImpl.GetByEmail"
		if err == sql.ErrNoRows {
			return nil, api.NewSingleMessageException(
				api.ENOTFOUND,
				op,
				"User Not Found",
				err,
			)
		}
		return nil, api.NewExceptionWithSourceLocation(op, "r.SQL.QueryRowContext", err)
	}

	return user, nil
}
