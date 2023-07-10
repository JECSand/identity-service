package repositories

import (
	"context"
	"github.com/JECSand/identity-service/command_service/config"
	"github.com/JECSand/identity-service/command_service/identity/models"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	blacklistQuery = `INSERT INTO blacklists (id, token, created_at) 
	VALUES ($1, $2, now()) RETURNING id, token, created_at`

	checkBlacklistQuery = `SELECT p.id, p.token, p.created_at,
	FROM blacklists p WHERE p.token = $1`

	countBlacklistQuery = `SELECT COUNT(*) from blacklists`
)

type blacklistRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

// NewBlacklistRepository ...
func NewBlacklistRepository(log logging.Logger, cfg *config.Config, db *pgxpool.Pool) *blacklistRepository {
	return &blacklistRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

// Count ...
func (p *blacklistRepository) Count(ctx context.Context) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistRepository.Count")
	defer span.Finish()
	type countRes struct {
		count int
	}
	var counted countRes
	if err := p.db.QueryRow(ctx, countBlacklistQuery).Scan(
		&counted.count,
	); err != nil {
		return 0, errors.Wrap(err, "Scan")
	}
	return counted.count, nil
}

// Create ...
func (p *blacklistRepository) Create(ctx context.Context, bl *models.Blacklist) (*models.Blacklist, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistRepository.Create")
	defer span.Finish()
	var created models.Blacklist
	if err := p.db.QueryRow(ctx, blacklistQuery, &bl.ID, &bl.AccessToken).Scan(
		&created.ID,
		&created.AccessToken,
		&created.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}
	return &created, nil
}

// GetByAccessToken ...
func (p *blacklistRepository) GetByAccessToken(ctx context.Context, accessToken string) (*models.Blacklist, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistRepository.GetByAccessToken")
	defer span.Finish()
	var found models.Blacklist
	if err := p.db.QueryRow(ctx, checkBlacklistQuery, accessToken).Scan(
		&found.ID,
		&found.AccessToken,
		&found.CreatedAt,
		&found.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &found, nil
}
