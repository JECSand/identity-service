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
	createMembershipQuery = `INSERT INTO memberships (id, user_id, group_id, status, member_role, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, now(), now()) RETURNING id, user_id, group_id, status, member_role, created_at, updated_at`

	updateMembershipQuery = `UPDATE memberships p SET 
                      status=COALESCE(NULLIF($2, 0), status), 
                      member_role=COALESCE(NULLIF($3, 0), member_role), 
                      updated_at = now()
                      WHERE id=$1
                      RETURNING id, user_id, group_id, status, member_role, created_at, updated_at`

	getMembershipByIdQuery = `SELECT p.id, p.user_id, p.group_id, p.status, p.member_role, p.created_at, p.updated_at 
	FROM memberships p WHERE p.id = $1`

	deleteMembershipByIdQuery = `DELETE FROM memberships WHERE id = $1`

	countMembershipsQuery = `SELECT COUNT(*) from memberships`

	getUserMembershipByIdQuery = `SELECT 
    	gen_random_uuid() AS id,
    	p.group_id,	
    	p.user_id, 
    	p.id AS membership_id, 
    	u.email,
    	u.username,
    	p.status, 
    	p.member_role AS role, 
    	p.created_at, 
    	p.updated_at 
	FROM memberships p 
	INNER JOIN users u ON p.user_id = u.id
	WHERE p.id = $1`

	getGroupMembershipByIdQuery = `SELECT 
    	gen_random_uuid() AS id,
    	p.user_id, 
    	p.group_id,	
    	p.id AS membership_id, 
    	g.group_name AS name,
    	g.description,
    	p.status, 
    	p.member_role AS role,
    	(p.user_id = g.creator_id) AS creator,
    	p.created_at, 
    	p.updated_at 
	FROM memberships p 
	INNER JOIN user_groups g ON p.group_id = g.id
	WHERE p.id = $1`
)

type membershipRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

// NewMembershipRepository ...
func NewMembershipRepository(log logging.Logger, cfg *config.Config, db *pgxpool.Pool) *membershipRepository {
	return &membershipRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

// Count ...
func (p *membershipRepository) Count(ctx context.Context) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.CountMemberships")
	defer span.Finish()
	type countRes struct {
		count int
	}
	var counted countRes
	if err := p.db.QueryRow(ctx, countMembershipsQuery).Scan(
		&counted.count,
	); err != nil {
		return 0, errors.Wrap(err, "Scan")
	}
	return counted.count, nil
}

// Create ...
func (p *membershipRepository) Create(ctx context.Context, membership *models.Membership) (*models.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.CreateMembership")
	defer span.Finish()
	var created models.Membership
	if err := p.db.QueryRow(ctx, createMembershipQuery, &membership.ID, &membership.UserID, &membership.GroupID, &membership.Status, membership.Role).Scan(
		&created.ID,
		&created.UserID,
		&created.GroupID,
		&created.Status,
		&created.Role,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}
	return &created, nil
}

// Update ...
func (p *membershipRepository) Update(ctx context.Context, membership *models.Membership) (*models.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.UpdateMembership")
	defer span.Finish()
	var updated models.Membership
	if err := p.db.QueryRow(
		ctx,
		updateMembershipQuery,
		&membership.ID,
		&membership.Status,
		&membership.Role,
	).Scan(&updated.ID, &updated.UserID, &updated.GroupID, &updated.Status, &updated.Role, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &updated, nil
}

// GetById ...
func (p *membershipRepository) GetById(ctx context.Context, uuid uuid.UUID) (*models.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.GetMembershipById")
	defer span.Finish()
	var found models.Membership
	if err := p.db.QueryRow(ctx, getMembershipByIdQuery, uuid).Scan(
		&found.ID,
		&found.UserID,
		&found.GroupID,
		&found.Status,
		&found.Role,
		&found.CreatedAt,
		&found.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &found, nil
}

// GetUserMembershipById ...
func (p *membershipRepository) GetUserMembershipById(ctx context.Context, uuid uuid.UUID) (*models.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.GetUserMembershipById")
	defer span.Finish()
	var found models.UserMembership
	if err := p.db.QueryRow(ctx, getUserMembershipByIdQuery, uuid).Scan(
		&found.ID,
		&found.GroupID,
		&found.UserID,
		&found.MembershipID,
		&found.Email,
		&found.Username,
		&found.Status,
		&found.Role,
		&found.CreatedAt,
		&found.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &found, nil
}

// GetGroupMembershipById ...
func (p *membershipRepository) GetGroupMembershipById(ctx context.Context, uuid uuid.UUID) (*models.GroupMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.GetGroupMembershipById")
	defer span.Finish()
	var found models.GroupMembership
	if err := p.db.QueryRow(ctx, getGroupMembershipByIdQuery, uuid).Scan(
		&found.ID,
		&found.UserID,
		&found.GroupID,
		&found.MembershipID,
		&found.Name,
		&found.Description,
		&found.Status,
		&found.Role,
		&found.Creator,
		&found.CreatedAt,
		&found.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}
	return &found, nil
}

// DeleteByID ...
func (p *membershipRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.DeleteMembershipByID")
	defer span.Finish()
	_, err := p.db.Exec(ctx, deleteMembershipByIdQuery, id)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}
	return nil
}
