package data

import (
	"context"
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

// groupEntity structures a user BSON document to save in a groups collection
type groupEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" validate:"required,min=3,max=500"`
	Description string             `bson:"description,omitempty" validate:"required"`
	CreatorID   primitive.ObjectID `bson:"creator_id,omitempty" validate:"required,min=3,max=250"`
	Active      bool               `bson:"active,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

// getID returns the unique identifier of the groupEntity
func (u *groupEntity) getID() primitive.ObjectID {
	return u.ID
}

// newGroupEntity initializes a new pointer to a groupEntity struct from a *models.Group struct
func newGroupEntity(u *entities.Group) (um *groupEntity, err error) {
	um = &groupEntity{
		Name:        u.Name,
		Description: u.Description,
		Active:      u.Active,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if utilities.CheckID(u.CreatorID) == nil {
		um.CreatorID, err = utilities.LoadObjectIDString(u.CreatorID)
	}
	if utilities.CheckID(u.ID) == nil {
		um.ID, err = utilities.LoadObjectIDString(u.ID)
	}
	return
}

// toRoot creates and return a new pointer to a models.Group JSON struct from a pointer to a BSON groupEntity
func (u *groupEntity) toRoot() *entities.Group {
	um := &entities.Group{
		Name:        u.Name,
		Description: u.Description,
		Active:      u.Active,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if utilities.CheckID(u.CreatorID.Hex()) == nil {
		um.CreatorID = utilities.LoadUUIDString(u.CreatorID)
	}
	if utilities.CheckID(u.ID.Hex()) == nil {
		um.ID = utilities.LoadUUIDString(u.ID)
	}
	return um
}

type groupRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewGroupRepository(log logging.Logger, cfg *config.Config, db *mongo.Client) *groupRepository {
	return &groupRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (p *groupRepository) Create(ctx context.Context, model *entities.Group) (*entities.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.CreateGroup")
	defer span.Finish()
	p.log.Info(model.ID)
	ent, err := newGroupEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Group{}, errors.Wrap(err, "newGroupEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Groups)
	_, err = collection.InsertOne(ctx, ent, &options.InsertOneOptions{})
	if err != nil {
		p.traceErr(span, err)
		return &entities.Group{}, errors.Wrap(err, "InsertOne")
	}
	return model, nil
}

func (p *groupRepository) Update(ctx context.Context, model *entities.Group) (*entities.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.UpdateGroup")
	defer span.Finish()
	ent, err := newGroupEntity(model)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Group{}, errors.Wrap(err, "newGroupEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Groups)
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)
	var updated entities.Group
	if err = collection.FindOneAndUpdate(ctx, bson.M{"_id": ent.ID}, bson.M{"$set": ent}, ops).Decode(&updated); err != nil {
		p.traceErr(span, err)
		return &entities.Group{}, errors.Wrap(err, "Decode")
	}
	return &updated, nil
}

func (p *groupRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.GetGroupById")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Groups)
	var ent groupEntity
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Group{}, errors.Wrap(err, "LoadObjectIDString")
	}
	if err = collection.FindOne(ctx, bson.M{"_id": oId}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.Group{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *groupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.DeleteGroup")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Groups)
	return collection.FindOneAndDelete(ctx, bson.M{"_id": oId}).Err()
}

func (p *groupRepository) Search(ctx context.Context, search string, pagination *utilities.Pagination) (*entities.GroupsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "groupRepository.Search")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Groups)
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
			bson.D{{Key: "description", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
		}},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupsList{}, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &entities.GroupsList{Groups: make([]*entities.Group, 0)}, nil
	}
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		p.traceErr(span, err)
		return &entities.GroupsList{}, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errCheck
	groups := make([]*entities.Group, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var u groupEntity
		if err = cursor.Decode(&u); err != nil {
			p.traceErr(span, err)
			return &entities.GroupsList{}, errors.Wrap(err, "Find")
		}
		groups = append(groups, u.toRoot())
	}
	if err = cursor.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return &entities.GroupsList{}, errors.Wrap(err, "cursor.Err")
	}
	return entities.NewGroupListWithPagination(groups, count, pagination), nil
}

func (p *groupRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
