package postgres

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db *database.PostgresDB
}

func NewUserRepository(db *database.PostgresDB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, bio, avatar, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash,
		user.Bio, user.Avatar, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, bio, avatar, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	if len(ids) == 0 {
		return []*models.User{}, nil
	}

	query := `
		SELECT id, name, email, password_hash, bio, avatar, created_at, updated_at
		FROM users
		WHERE id = ANY($1)
	`

	rows, err := r.db.Pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, bio, avatar, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET name = $2, email = $3, password_hash = $4, bio = $5, avatar = $6, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash, user.Bio, user.Avatar,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, bio, avatar, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetAll(ctx context.Context, limit int, after *string) ([]*models.User, error) {
	var query string
	var rows pgx.Rows
	var err error

	if after != nil && *after != "" {
		query = `
			SELECT id, name, email, password_hash, bio, avatar, created_at, updated_at
			FROM users
			WHERE id > $1
			ORDER BY id ASC
			LIMIT $2
		`
		afterID, parseErr := uuid.Parse(*after)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid cursor: %w", parseErr)
		}
		rows, err = r.db.Pool.Query(ctx, query, afterID, limit)
	} else {
		query = `
			SELECT id, name, email, password_hash, bio, avatar, created_at, updated_at
			FROM users
			ORDER BY id ASC
			LIMIT $1
		`
		rows, err = r.db.Pool.Query(ctx, query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&user.Bio, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
