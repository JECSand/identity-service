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

// groupMembershipEntity structures a user BSON document to save in a groupMemberships aggregate collection
type groupMembershipEntity struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	UserID       primitive.ObjectID     `bson:"user_id,omitempty" validate:"required"`
	GroupID      primitive.ObjectID     `bson:"group_id,omitempty" validate:"required"`
	MembershipID primitive.ObjectID     `bson:"membership_id,omitempty" validate:"required"`
	Name         string                 `bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description  string                 `bson:"description,omitempty" validate:"required,min=3,max=250"`
	Status       enums.MembershipStatus `bson:"status,omitempty"`
	Role         enums.Role             `bson:"role,omitempty"`
	Creator      bool                   `bson:"creator,omitempty"`
	CreatedAt    time.Time              `bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `bson:"updated_at,omitempty"`
}

// getID returns the unique identifier of the groupMembershipEntity
func (u *groupMembershipEntity) getID() primitive.ObjectID {
	return u.ID
}

// bsonFilter generates a bson filter for MongoDB queries from the groupMembershipEntity data
func (u *groupMembershipEntity) bsonFilter() (doc bson.D, err error) {
	if utilities.CheckObjectID(u.ID) == nil {
		doc = bson.D{{"_id", u.ID}}
	} else if utilities.CheckObjectID(u.UserID) == nil {
		doc = bson.D{{"user_id", u.UserID}}
	} else if utilities.CheckObjectID(u.GroupID) == nil {
		doc = bson.D{{"group_id", u.GroupID}}
	} else if utilities.CheckObjectID(u.MembershipID) == nil {
		doc = bson.D{{"membership_id", u.MembershipID}}
	}
	return
}

// newGroupMembershipEntity initializes a new pointer to an groupMembershipEntity struct from a models.GroupMembership
func newGroupMembershipEntity(u *entities.GroupMembership) (um *groupMembershipEntity, err error) {
	um = &groupMembershipEntity{
		Name:        u.Name,
		Description: u.Description,
		Status:      u.Status,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if utilities.CheckID(u.UserID) == nil {
		um.UserID, err = utilities.LoadObjectIDString(u.UserID)
	}
	if utilities.CheckID(u.GroupID) == nil {
		um.GroupID, err = utilities.LoadObjectIDString(u.GroupID)
	}
	if utilities.CheckID(u.MembershipID) == nil {
		um.MembershipID, err = utilities.LoadObjectIDString(u.MembershipID)
	}
	if utilities.CheckID(u.ID) == nil {
		um.ID, err = utilities.LoadObjectIDString(u.ID)
	}
	return
}

// toRoot creates and returns a new pointer to a models.GroupMembership JSON struct from a groupMembershipEntity
func (u *groupMembershipEntity) toRoot() *entities.GroupMembership {
	um := &entities.GroupMembership{
		Name:        u.Name,
		Description: u.Description,
		Status:      u.Status,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if utilities.CheckID(u.UserID.Hex()) == nil {
		um.UserID = utilities.LoadUUIDString(u.UserID)
	}
	if utilities.CheckID(u.GroupID.Hex()) == nil {
		um.UserID = utilities.LoadUUIDString(u.GroupID)
	}
	if utilities.CheckID(u.MembershipID.Hex()) == nil {
		um.MembershipID = utilities.LoadUUIDString(u.MembershipID)
	}
	if utilities.CheckID(u.ID.Hex()) == nil {
		um.ID = utilities.LoadUUIDString(u.ID)
	}
	return um
}

type groupMembershipRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewGroupMembershipRepository(log logging.Logger, cfg *config.Config, db *mongo.Client) *groupMembershipRepository {
	return &groupMembershipRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (p *groupMembershipRepository) Create(ctx context.Context, model *entities.GroupMembership) (*entities.GroupMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.CreateGroupMembership")
	defer span.Finish()
	p.log.Info(model.ID)
	ent, err := newGroupMembershipEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembership{}, errors.Wrap(err, "newGroupMembershipEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	_, err = collection.InsertOne(ctx, ent, &options.InsertOneOptions{})
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembership{}, errors.Wrap(err, "InsertOne")
	}
	return model, nil
}

func (p *groupMembershipRepository) UpdateMany(ctx context.Context, filter *entities.GroupMembership, update *entities.GroupMembership) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.UpdateGroupMemberships")
	defer span.Finish()
	up, err := newGroupMembershipEntity(update)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newGroupMembershipEntity")
	}
	f, err := newGroupMembershipEntity(filter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newGroupMembershipEntity")
	}
	bsonFilter, err := f.bsonFilter()
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "bsonFilter")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	_, err = collection.UpdateMany(ctx, bsonFilter, bson.M{"$set": up})
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "collection.UpdateMany")
	}
	return nil
}

func (p *groupMembershipRepository) Update(ctx context.Context, model *entities.GroupMembership) (*entities.GroupMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.UpdateGroupMembership")
	defer span.Finish()
	ent, err := newGroupMembershipEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembership{}, errors.Wrap(err, "newGroupMembershipEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)
	var updated entities.GroupMembership
	if err = collection.FindOneAndUpdate(ctx, bson.M{"_id": ent.ID}, bson.M{"$set": ent}, ops).Decode(&updated); err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembership{}, errors.Wrap(err, "Decode")
	}
	return &updated, nil
}

func (p *groupMembershipRepository) GetById(ctx context.Context, id uuid.UUID, idType enums.ReadTableIdType) (*entities.GroupMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.GetGroupMembershipById")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	var ent groupMembershipEntity
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembership{}, errors.Wrap(err, "LoadObjectIDString")
	}
	if err = collection.FindOne(ctx, bson.M{idType.KeyString(): oId}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembership{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *groupMembershipRepository) GetByUserId(ctx context.Context, userId uuid.UUID, pagination *utilities.Pagination) (*entities.GroupMembershipsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.GetByUserId")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	filter := bson.M{"user_id": userId}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembershipsList{}, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &entities.GroupMembershipsList{GroupMemberships: make([]*entities.GroupMembership, 0)}, nil
	}
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembershipsList{}, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errCheck
	groups := make([]*entities.GroupMembership, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var u groupMembershipEntity
		if err = cursor.Decode(&u); err != nil {
			p.traceErr(span, err)
			return &entities.GroupMembershipsList{}, errors.Wrap(err, "Find")
		}
		groups = append(groups, u.toRoot())
	}
	if err = cursor.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return &entities.GroupMembershipsList{}, errors.Wrap(err, "cursor.Err")
	}
	return entities.NewGroupMembershipListWithPagination(groups, count, pagination), nil
}

func (p *groupMembershipRepository) GetByGroupId(ctx context.Context, groupId uuid.UUID, pagination *utilities.Pagination) (*entities.GroupMembershipsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.GetByGroupId")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	filter := bson.M{"group_id": groupId}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembershipsList{}, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &entities.GroupMembershipsList{GroupMemberships: make([]*entities.GroupMembership, 0)}, nil
	}
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupMembershipsList{}, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errCheck
	groups := make([]*entities.GroupMembership, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var u groupMembershipEntity
		if err = cursor.Decode(&u); err != nil {
			p.traceErr(span, err)
			return &entities.GroupMembershipsList{}, errors.Wrap(err, "Find")
		}
		groups = append(groups, u.toRoot())
	}
	if err = cursor.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return &entities.GroupMembershipsList{}, errors.Wrap(err, "cursor.Err")
	}
	return entities.NewGroupMembershipListWithPagination(groups, count, pagination), nil
}

func (p *groupMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.DeleteGroupMembership")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	return collection.FindOneAndDelete(ctx, bson.M{"_id": oId}).Err()
}

func (p *groupMembershipRepository) DeleteByMembershipId(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.DeleteGroupMembershipByMembershipId")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	return collection.FindOneAndDelete(ctx, bson.M{"MembershipID": oId}).Err()
}

func (p *groupMembershipRepository) DeleteMany(ctx context.Context, filter *entities.GroupMembership) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupMembershipRepository.DeleteGroupMemberships")
	defer span.Finish()
	f, err := newGroupMembershipEntity(filter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newGroupMembershipEntity")
	}
	bsonFilter, err := f.bsonFilter()
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "bsonFilter")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.GroupMemberships)
	_, err = collection.DeleteMany(ctx, bsonFilter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "collection.DeleteMany")
	}
	return nil
}

func (p *groupMembershipRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
