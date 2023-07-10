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

// userMembershipEntity structures a user BSON document to save in a userMemberships aggregate collection
type userMembershipEntity struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	GroupID      primitive.ObjectID     `bson:"group_id,omitempty"`
	UserID       primitive.ObjectID     `bson:"user_id,omitempty"`
	MembershipID primitive.ObjectID     `bson:"membership_id,omitempty"`
	Email        string                 `bson:"email,omitempty" validate:"required,min=3,max=250"`
	Username     string                 `bson:"username,omitempty" validate:"required,min=3,max=250"`
	Status       enums.MembershipStatus `bson:"status,omitempty"`
	Role         enums.Role             `bson:"role,omitempty"`
	CreatedAt    time.Time              `bson:"created_at,omitempty"`
	UpdatedAt    time.Time              `bson:"updated_at,omitempty"`
}

// getID returns the unique identifier of the userMembershipEntity
func (u *userMembershipEntity) getID() primitive.ObjectID {
	return u.ID
}

// bsonFilter generates a bson filter for MongoDB queries from the userMembershipEntity data
func (u *userMembershipEntity) bsonFilter() (doc bson.D, err error) {
	if utilities.CheckObjectID(u.ID) == nil {
		doc = bson.D{{"_id", u.ID}}
	} else if utilities.CheckObjectID(u.GroupID) == nil {
		doc = bson.D{{"group_id", u.GroupID}}
	} else if utilities.CheckObjectID(u.UserID) == nil {
		doc = bson.D{{"user_id", u.UserID}}
	} else if utilities.CheckObjectID(u.MembershipID) == nil {
		doc = bson.D{{"membership_id", u.MembershipID}}
	}
	return
}

// newUserMembershipEntity initializes a new pointer to an userMembershipEntity struct from a models.UserMembership
func newUserMembershipEntity(u *entities.UserMembership) (um *userMembershipEntity, err error) {
	um = &userMembershipEntity{
		Email:     u.Email,
		Username:  u.Username,
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if utilities.CheckID(u.GroupID) == nil {
		um.GroupID, err = utilities.LoadObjectIDString(u.GroupID)
	}
	if utilities.CheckID(u.UserID) == nil {
		um.UserID, err = utilities.LoadObjectIDString(u.UserID)
	}
	if utilities.CheckID(u.MembershipID) == nil {
		um.MembershipID, err = utilities.LoadObjectIDString(u.MembershipID)
	}
	if utilities.CheckID(u.ID) == nil {
		um.ID, err = utilities.LoadObjectIDString(u.ID)
	}
	return
}

// toRoot creates and returns a new pointer to a models.UserMembership JSON struct from a userMembershipEntity
func (u *userMembershipEntity) toRoot() *entities.UserMembership {
	um := &entities.UserMembership{
		Email:     u.Email,
		Username:  u.Username,
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if utilities.CheckID(u.GroupID.Hex()) == nil {
		um.UserID = utilities.LoadUUIDString(u.GroupID)
	}
	if utilities.CheckID(u.UserID.Hex()) == nil {
		um.UserID = utilities.LoadUUIDString(u.UserID)
	}
	if utilities.CheckID(u.MembershipID.Hex()) == nil {
		um.MembershipID = utilities.LoadUUIDString(u.MembershipID)
	}
	if utilities.CheckID(u.ID.Hex()) == nil {
		um.ID = utilities.LoadUUIDString(u.ID)
	}
	return um
}

type userMembershipRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewUserMembershipRepository(log logging.Logger, cfg *config.Config, db *mongo.Client) *userMembershipRepository {
	return &userMembershipRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (p *userMembershipRepository) Create(ctx context.Context, model *entities.UserMembership) (*entities.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.CreateUserMembership")
	defer span.Finish()
	p.log.Info(model.ID)
	ent, err := newUserMembershipEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembership{}, errors.Wrap(err, "newUserMembershipEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	_, err = collection.InsertOne(ctx, ent, &options.InsertOneOptions{})
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembership{}, errors.Wrap(err, "InsertOne")
	}
	return model, nil
}

func (p *userMembershipRepository) UpdateMany(ctx context.Context, filter *entities.UserMembership, update *entities.UserMembership) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.UpdateUserMemberships")
	defer span.Finish()
	up, err := newUserMembershipEntity(update)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newUserMembershipEntity")
	}
	f, err := newUserMembershipEntity(filter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newUserMembershipEntity")
	}
	bsonFilter, err := f.bsonFilter()
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "bsonFilter")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	_, err = collection.UpdateMany(ctx, bsonFilter, bson.M{"$set": up})
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "collection.UpdateMany")
	}
	return nil
}

func (p *userMembershipRepository) Update(ctx context.Context, model *entities.UserMembership) (*entities.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.UpdateUserMembership")
	defer span.Finish()
	ent, err := newUserMembershipEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembership{}, errors.Wrap(err, "newUserMembershipEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)
	var updated entities.UserMembership
	if err = collection.FindOneAndUpdate(ctx, bson.M{"_id": ent.ID}, bson.M{"$set": ent}, ops).Decode(&updated); err != nil {
		p.traceErr(span, err)
		return &entities.UserMembership{}, errors.Wrap(err, "Decode")
	}
	return &updated, nil
}

func (p *userMembershipRepository) GetById(ctx context.Context, id uuid.UUID, idType enums.ReadTableIdType) (*entities.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.GetUserMembershipById")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	var ent userMembershipEntity
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembership{}, errors.Wrap(err, "LoadObjectIDString")
	}
	if err = collection.FindOne(ctx, bson.M{idType.KeyString(): oId}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.UserMembership{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *userMembershipRepository) GetByUserId(ctx context.Context, userId uuid.UUID, pagination *utilities.Pagination) (*entities.UserMembershipsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.GetByUserId")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	filter := bson.M{"user_id": userId}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembershipsList{}, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &entities.UserMembershipsList{UserMemberships: make([]*entities.UserMembership, 0)}, nil
	}
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembershipsList{}, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errCheck
	groups := make([]*entities.UserMembership, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var u userMembershipEntity
		if err = cursor.Decode(&u); err != nil {
			p.traceErr(span, err)
			return &entities.UserMembershipsList{}, errors.Wrap(err, "Find")
		}
		groups = append(groups, u.toRoot())
	}
	if err = cursor.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return &entities.UserMembershipsList{}, errors.Wrap(err, "cursor.Err")
	}
	return entities.NewUserMembershipListWithPagination(groups, count, pagination), nil
}

func (p *userMembershipRepository) GetByGroupId(ctx context.Context, groupId uuid.UUID, pagination *utilities.Pagination) (*entities.UserMembershipsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.GetByGroupId")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	filter := bson.M{"group_id": groupId}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembershipsList{}, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &entities.UserMembershipsList{UserMemberships: make([]*entities.UserMembership, 0)}, nil
	}
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		p.traceErr(span, err)
		return &entities.UserMembershipsList{}, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errCheck
	groups := make([]*entities.UserMembership, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var u userMembershipEntity
		if err = cursor.Decode(&u); err != nil {
			p.traceErr(span, err)
			return &entities.UserMembershipsList{}, errors.Wrap(err, "Find")
		}
		groups = append(groups, u.toRoot())
	}
	if err = cursor.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return &entities.UserMembershipsList{}, errors.Wrap(err, "cursor.Err")
	}
	return entities.NewUserMembershipListWithPagination(groups, count, pagination), nil
}

func (p *userMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.DeleteUserMembership")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	return collection.FindOneAndDelete(ctx, bson.M{"_id": oId}).Err()
}

func (p *userMembershipRepository) DeleteByMembershipId(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.DeleteUserMembershipByMembershipId")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	return collection.FindOneAndDelete(ctx, bson.M{"MembershipID": oId}).Err()
}

func (p *userMembershipRepository) DeleteMany(ctx context.Context, filter *entities.UserMembership) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userMembershipRepository.DeleteUserMemberships")
	defer span.Finish()
	f, err := newUserMembershipEntity(filter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "newUserMembershipEntity")
	}
	bsonFilter, err := f.bsonFilter()
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "bsonFilter")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.UserMemberships)
	_, err = collection.DeleteMany(ctx, bsonFilter)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "collection.DeleteMany")
	}
	return nil
}

func (p *userMembershipRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
