package data

import (
	"context"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// membershipEntity structures a user BSON document to save in a memberships collection
type membershipEntity struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty"`
	UserID    primitive.ObjectID     `bson:"user_id,omitempty" validate:"required"`
	GroupID   primitive.ObjectID     `bson:"group_id,omitempty" validate:"required"`
	Status    enums.MembershipStatus `bson:"status,omitempty"`
	Role      enums.Role             `bson:"role,omitempty"`
	CreatedAt time.Time              `bson:"created_at,omitempty"`
	UpdatedAt time.Time              `bson:"updated_at,omitempty"`
}

// getID returns the unique identifier of the membershipEntity
func (u *membershipEntity) getID() primitive.ObjectID {
	return u.ID
}

// bsonFilter generates a bson filter for MongoDB queries from the membershipEntity data
func (u *membershipEntity) bsonFilter() (doc bson.D, err error) {
	if utilities.CheckObjectID(u.ID) == nil {
		doc = bson.D{{"_id", u.ID}}
	} else if utilities.CheckObjectID(u.UserID) == nil {
		doc = bson.D{{"user_id", u.UserID}}
	} else if utilities.CheckObjectID(u.GroupID) == nil {
		doc = bson.D{{"group_id", u.GroupID}}
	}
	return
}

// newMembershipEntity initializes a new pointer to a membershipEntity struct from a models.Membership
func newMembershipEntity(u *entities.Membership) (um *membershipEntity, err error) {
	um = &membershipEntity{
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if utilities.CheckID(u.UserID) == nil {
		um.UserID, err = utilities.LoadObjectIDString(u.UserID)
	}
	if utilities.CheckID(u.GroupID) == nil {
		um.GroupID, err = utilities.LoadObjectIDString(u.GroupID)
	}
	if utilities.CheckID(u.ID) == nil {
		um.ID, err = utilities.LoadObjectIDString(u.ID)
	}
	return
}

// toRoot creates and return a new pointer to a *models.Membership JSON struct from a pointer to a BSON membershipEntity
func (u *membershipEntity) toRoot() *entities.Membership {
	um := &entities.Membership{
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if utilities.CheckID(u.UserID.Hex()) == nil {
		um.UserID = utilities.LoadUUIDString(u.UserID)
	}
	if utilities.CheckID(u.GroupID.Hex()) == nil {
		um.UserID = utilities.LoadUUIDString(u.GroupID)
	}
	if utilities.CheckID(u.ID.Hex()) == nil {
		um.ID = utilities.LoadUUIDString(u.ID)
	}
	return um
}

type membershipRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewMembershipRepository(log logging.Logger, cfg *config.Config, db *mongo.Client) *membershipRepository {
	return &membershipRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (p *membershipRepository) Create(ctx context.Context, model *entities.Membership) (*entities.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.CreateMembership")
	defer span.Finish()
	p.log.Info(model.ID)
	ent, err := newMembershipEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Membership{}, errors.Wrap(err, "newMembershipEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Memberships)
	_, err = collection.InsertOne(ctx, ent, &options.InsertOneOptions{})
	if err != nil {
		p.traceErr(span, err)
		return &entities.Membership{}, errors.Wrap(err, "InsertOne")
	}
	return model, nil
}

func (p *membershipRepository) UpdateMany(ctx context.Context, filter *entities.Membership, update *entities.Membership) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.UpdateMemberships")
	defer span.Finish()
	up, err := newMembershipEntity(update)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newMembershipEntity")
	}
	f, err := newMembershipEntity(filter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newMembershipEntity")
	}
	bsonFilter, err := f.bsonFilter()
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "bsonFilter")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Memberships)
	_, err = collection.UpdateMany(ctx, bsonFilter, bson.M{"$set": up})
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "collection.UpdateMany")
	}
	return nil
}

func (p *membershipRepository) Update(ctx context.Context, model *entities.Membership) (*entities.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.UpdateMembership")
	defer span.Finish()
	ent, err := newMembershipEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Membership{}, errors.Wrap(err, "newMembershipEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Memberships)
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)
	var updated entities.Membership
	if err = collection.FindOneAndUpdate(ctx, bson.M{"_id": ent.ID}, bson.M{"$set": ent}, ops).Decode(&updated); err != nil {
		p.traceErr(span, err)
		return &entities.Membership{}, errors.Wrap(err, "Decode")
	}
	return &updated, nil
}

func (p *membershipRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.GetMembershipById")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Memberships)
	var ent membershipEntity
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Membership{}, errors.Wrap(err, "LoadObjectIDString")
	}
	if err = collection.FindOne(ctx, bson.M{"_id": oId}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.Membership{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *membershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.DeleteMembership")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Memberships)
	return collection.FindOneAndDelete(ctx, bson.M{"_id": oId}).Err()
}

func (p *membershipRepository) DeleteMany(ctx context.Context, filter *entities.Membership) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "membershipRepository.DeleteMemberships")
	defer span.Finish()
	f, err := newMembershipEntity(filter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newMembershipEntity")
	}
	bsonFilter, err := f.bsonFilter()
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "bsonFilter")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Memberships)
	_, err = collection.DeleteMany(ctx, bsonFilter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "collection.DeleteMany")
	}
	return nil
}

func (p *membershipRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
