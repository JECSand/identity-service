package repositories

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	createUserQuery = `INSERT INTO users (id, email, username, password, root, active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, now(), now()) RETURNING id, email, username, password, root, active, created_at, updated_at`

	updateUserQuery = `UPDATE users p SET 
                      email=COALESCE(NULLIF($2, ''), email), 
                      username=COALESCE(NULLIF($3, ''), username), 
                      updated_at = now()
                      WHERE id=$1
                      RETURNING id, email, username, root, active, created_at, updated_at`

	updateUserPasswordQuery = `UPDATE users p SET 
                      password=COALESCE(NULLIF($2, ''), password), 
                      updated_at = now()
                      WHERE id=$1
                      RETURNING id, email, username, root, active, created_at, updated_at`

	getUserByIdQuery = `SELECT p.id, p.email, p.username, p.password, p.root, p.active, p.created_at, p.updated_at 
	FROM users p WHERE p.id = $1`

	deleteUserByIdQuery = `DELETE FROM users WHERE id = $1`

	countUsersQuery = `SELECT COUNT(*) from users`
)

type userRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

// NewUserRepository ...
func NewUserRepository(log logging.Logger, cfg *config.Config, db *pgxpool.Pool) *userRepository {
	return &userRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

// Count ...
func (p *userRepository) Count(ctx context.Context) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.CountUsers")
	defer span.Finish()
	type countRes struct {
		count int
	}
	var counted countRes
	if err := p.db.QueryRow(ctx, countUsersQuery).Scan(
		&counted.count,
	); err != nil {
		return 0, errors.Wrap(err, "Scan")
	}
	return counted.count, nil
}

// Create ...
func (p *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.CreateUser")
	defer span.Finish()
	var created models.User
	if err := p.db.QueryRow(ctx, createUserQuery, &user.ID, &user.Email, &user.Username, &user.Password, user.Root, user.Active).Scan(
		&created.ID,
		&created.Email,
		&created.Username,
		&created.Password,
		&created.Root,
		&created.Active,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}
	return &created, nil
}

// Update ...
func (p *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.UpdateUser")
	defer span.Finish()
	var updated models.User
	if err := p.db.QueryRow(
		ctx,
		updateUserQuery,
		&user.ID,
		&user.Email,
		&user.Username,
	).Scan(&updated.ID, &updated.Email, &updated.Username, &updated.Root, &updated.Active, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &updated, nil
}

// UpdatePassword ...
func (p *userRepository) UpdatePassword(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.UpdatePassword")
	defer span.Finish()
	var updated models.User
	if err := p.db.QueryRow(
		ctx,
		updateUserPasswordQuery,
		&user.ID,
		&user.Password,
	).Scan(&updated.ID, &updated.Email, &updated.Username, &updated.Root, &updated.Active, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &updated, nil
}

// GetById ...
func (p *userRepository) GetById(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.GetUserById")
	defer span.Finish()
	var found models.User
	if err := p.db.QueryRow(ctx, getUserByIdQuery, uuid).Scan(
		&found.ID,
		&found.Email,
		&found.Username,
		&found.Password,
		&found.Root,
		&found.Active,
		&found.CreatedAt,
		&found.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &found, nil
}

// DeleteByID ...
func (p *userRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.DeleteUserByID")
	defer span.Finish()
	_, err := p.db.Exec(ctx, deleteUserByIdQuery, id)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}
	return nil
}
