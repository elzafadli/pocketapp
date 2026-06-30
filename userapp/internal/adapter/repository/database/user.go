package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"userapp/internal/adapter/repository"
	"userapp/internal/adapter/repository/database/querier"
	"userapp/internal/domain/user"

	"github.com/google/uuid"
	"github.com/runsystemid/golog"
)

type userDB struct {
	ID              string     `db:"id"`
	Name            string     `db:"name"`
	TenantDefault   string     `db:"tenant_default"`
	ActiveIndicator string     `db:"active_indicator"`
	Email           string     `db:"email"`
	Password        string     `db:"password"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

type UserRepository struct {
	DB      repository.Sqlx     `inject:"database"`
	Querier querier.UserQuerier `inject:"userQuerier"`
}

func (r *UserRepository) Startup() error {
	return nil
}

func (r *UserRepository) Shutdown() error {
	return nil
}

func (r *UserRepository) Create(ctx context.Context, data *user.User) (*user.User, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	query, args, err := r.Querier.Create(data)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query create user: %s", err.Error()), err)
		return nil, user.ErrDataCreate
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		if r.DB.IsErrorDuplicate(err) {
			return nil, user.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error create user: %s", err.Error()), err)
		return nil, user.ErrDataCreate
	}

	return data, nil
}

func (r *UserRepository) Update(ctx context.Context, data *user.User) error {
	query, args, err := r.Querier.Update(data)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query update user: %s", err.Error()), err)
		return user.ErrDataUpdate
	}

	result, err := r.DB.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		if r.DB.IsErrorDuplicate(err) {
			return user.ErrDataAlreadyExists
		}
		golog.Error(ctx, fmt.Sprintf("Error update user: %s", err.Error()), err)
		return user.ErrDataUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error getting rows affected in update user: %s", err.Error()), err)
		return user.ErrDataUpdate
	}

	if rowsAffected == 0 {
		return user.ErrDataNotFound
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query, args, err := r.Querier.GetByID(id)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query get user by ID: %s", err.Error()), err)
		return nil, user.ErrDataGet
	}

	var dbUser userDB
	err = r.DB.GetDB().GetContext(ctx, &dbUser, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrDataNotFound
		}
		golog.Error(ctx, fmt.Sprintf("Error get user by ID: %s", err.Error()), err)
		return nil, user.ErrDataGet
	}

	return mapToUserDomain(&dbUser), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query, args, err := r.Querier.GetByEmail(email)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query get user by email: %s", err.Error()), err)
		return nil, user.ErrDataGet
	}

	var dbUser userDB
	err = r.DB.GetDB().GetContext(ctx, &dbUser, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrDataNotFound
		}
		golog.Error(ctx, fmt.Sprintf("Error get user by email: %s", err.Error()), err)
		return nil, user.ErrDataGet
	}

	return mapToUserDomain(&dbUser), nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.Querier.Delete(id)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query delete user: %s", err.Error()), err)
		return user.ErrDataDelete
	}

	result, err := r.DB.GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error delete user: %s", err.Error()), err)
		return user.ErrDataDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error getting rows affected in delete user: %s", err.Error()), err)
		return user.ErrDataDelete
	}

	if rowsAffected == 0 {
		return user.ErrDataNotFound
	}

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context, filter map[string]interface{}) ([]*user.User, error) {
	query, args, err := r.Querier.List(filter)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query get all users: %s", err.Error()), err)
		return nil, user.ErrDataGet
	}

	dbUsers := make([]*userDB, 0)
	err = r.DB.GetDB().SelectContext(ctx, &dbUsers, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error get all users: %s", err.Error()), err)
		return nil, user.ErrDataGet
	}

	res := make([]*user.User, len(dbUsers))
	for i, dbU := range dbUsers {
		res[i] = mapToUserDomain(dbU)
	}

	return res, nil
}

func (r *UserRepository) Count(ctx context.Context, filter map[string]interface{}) (uint64, error) {
	query, args, err := r.Querier.Count(filter)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error build query count users: %s", err.Error()), err)
		return 0, user.ErrDataGet
	}

	var count uint64
	err = r.DB.GetDB().GetContext(ctx, &count, query, args...)
	if err != nil {
		golog.Error(ctx, fmt.Sprintf("Error count users: %s", err.Error()), err)
		return 0, user.ErrDataGet
	}

	return count, nil
}

func mapToUserDomain(dbUser *userDB) *user.User {
	return &user.User{
		ID:              dbUser.ID,
		Name:            dbUser.Name,
		TenantDefault:   dbUser.TenantDefault,
		ActiveIndicator: dbUser.ActiveIndicator,
		Email:           dbUser.Email,
		Password:        dbUser.Password,
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       dbUser.UpdatedAt,
		DeletedAt:       dbUser.DeletedAt,
	}
}
