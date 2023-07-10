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
	createGroupQuery = `INSERT INTO user_groups (id, group_name, description, creator_id, active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, now(), now()) RETURNING id, group_name, description, creator_id, active, created_at, updated_at`

	updateGroupQuery = `UPDATE user_groups p SET 
                      group_name=COALESCE(NULLIF($2, ''), group_name), 
                      description=COALESCE(NULLIF($3, ''), description), 
                      updated_at = now()
                      WHERE id=$1
                      RETURNING id, group_name, description, creator_id, active, created_at, updated_at`

	getGroupByIdQuery = `SELECT p.id, p.group_name AS name, p.description, p.creator_id, p.active, p.created_at, p.updated_at 
	FROM user_groups p WHERE p.id = $1`

	deleteGroupByIdQuery = `DELETE FROM user_groups WHERE id = $1`

	countGroupsQuery = `SELECT COUNT(*) from user_groups`
)

type groupRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

// NewGroupRepository ...
func NewGroupRepository(log logging.Logger, cfg *config.Config, db *pgxpool.Pool) *groupRepository {
	return &groupRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

// Count ...
func (p *groupRepository) Count(ctx context.Context) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.CountGroups")
	defer span.Finish()
	type countRes struct {
		count int
	}
	var counted countRes
	if err := p.db.QueryRow(ctx, countGroupsQuery).Scan(
		&counted.count,
	); err != nil {
		return 0, errors.Wrap(err, "Scan")
	}
	return counted.count, nil
}

// Create ...
func (p *groupRepository) Create(ctx context.Context, group *models.Group) (*models.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.CreateGroup")
	defer span.Finish()
	var created models.Group
	if err := p.db.QueryRow(ctx, createGroupQuery, &group.ID, &group.Name, &group.Description, &group.CreatorID, group.Active).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.CreatorID,
		&created.Active,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}
	return &created, nil
}

// Update ...
func (p *groupRepository) Update(ctx context.Context, group *models.Group) (*models.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.UpdateGroup")
	defer span.Finish()
	var updated models.Group
	if err := p.db.QueryRow(
		ctx,
		updateGroupQuery,
		&group.ID,
		&group.Name,
		&group.Description,
	).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.CreatorID, &updated.Active, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &updated, nil
}

// GetById ...
func (p *groupRepository) GetById(ctx context.Context, uuid uuid.UUID) (*models.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.GetGroupById")
	defer span.Finish()
	var found models.Group
	if err := p.db.QueryRow(ctx, getGroupByIdQuery, uuid).Scan(
		&found.ID,
		&found.Name,
		&found.Description,
		&found.CreatorID,
		&found.Active,
		&found.CreatedAt,
		&found.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &found, nil
}

// DeleteByID ...
func (p *groupRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.DeleteGroupByID")
	defer span.Finish()
	_, err := p.db.Exec(ctx, deleteGroupByIdQuery, id)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}
	return nil
}
