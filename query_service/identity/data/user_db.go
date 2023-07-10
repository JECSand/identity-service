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

// userEntity structures a user BSON document to save in a users collection
type userEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username,omitempty" validate:"required,min=3,max=500"`
	Password  string             `bson:"password,omitempty" validate:"required"`
	Email     string             `bson:"email,omitempty" validate:"required,min=3,max=250"`
	Root      bool               `bson:"root,omitempty"`
	Active    bool               `bson:"active,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

// getID returns the unique identifier of the userEntity
func (u *userEntity) getID() primitive.ObjectID {
	return u.ID
}

// newUserEntity initializes a new pointer to a userEntity struct from a pointer to a JSON models.User struct
func newUserEntity(u *entities.User) (um *userEntity, err error) {
	um = &userEntity{
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Root:      u.Root,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if utilities.CheckID(u.ID) == nil {
		um.ID, err = utilities.LoadObjectIDString(u.ID)
	}
	return
}

// toRoot creates and return a new pointer to a models.User JSON struct from a pointer to a BSON userEntity
func (u *userEntity) toRoot() *entities.User {
	um := &entities.User{
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		Root:      u.Root,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if utilities.CheckID(u.ID.Hex()) == nil {
		um.ID = utilities.LoadUUIDString(u.ID)
	}
	return um
}

func loadUserEntities(ms []*userEntity) (users []*entities.User) {
	for _, m := range ms {
		users = append(users, m.toRoot())
	}
	return
}

type userRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewUserRepository(log logging.Logger, cfg *config.Config, db *mongo.Client) *userRepository {
	return &userRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (p *userRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.CreateUser")
	defer span.Finish()
	p.log.Info(user.ID)
	//tID := utilities.
	ent, err := newUserEntity(user)
	if err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "newUserEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Users)
	_, err = collection.InsertOne(ctx, ent, &options.InsertOneOptions{})
	if err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "InsertOne")
	}
	return user, nil
}

func (p *userRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.UpdateUser")
	defer span.Finish()
	ent, err := newUserEntity(user)
	if err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "newUserEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Users)
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)
	var updated entities.User
	if err = collection.FindOneAndUpdate(ctx, bson.M{"_id": ent.ID}, bson.M{"$set": ent}, ops).Decode(&updated); err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "Decode")
	}
	return &updated, nil
}

func (p *userRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.GetUserById")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Users)
	var ent userEntity
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "LoadObjectIDString")
	}
	if err = collection.FindOne(ctx, bson.M{"_id": oId}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.GetByEmail")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Users)
	var ent userEntity
	if err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.DeleteUser")
	defer span.Finish()
	oId, err := utilities.LoadObjectID(id)
	if err != nil {
		p.traceErr(span, err)
		return errors.Wrap(err, "LoadObjectIDString")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Users)
	return collection.FindOneAndDelete(ctx, bson.M{"_id": oId}).Err()
}

func (p *userRepository) Search(ctx context.Context, search string, pagination *utilities.Pagination) (*entities.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.Search")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Users)
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
			bson.D{{Key: "description", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
		}},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		p.traceErr(span, err)
		return &entities.UsersList{}, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &entities.UsersList{Users: make([]*entities.User, 0)}, nil
	}
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		p.traceErr(span, err)
		return &entities.UsersList{}, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx) // nolint: errCheck
	users := make([]*entities.User, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var u userEntity
		if err = cursor.Decode(&u); err != nil {
			p.traceErr(span, err)
			return &entities.UsersList{}, errors.Wrap(err, "Find")
		}
		users = append(users, u.toRoot())
	}
	if err = cursor.Err(); err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
		return &entities.UsersList{}, errors.Wrap(err, "cursor.Err")
	}
	return entities.NewUserListWithPagination(users, count, pagination), nil
}

func (p *userRepository) Authenticate(ctx context.Context, email string, password string) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRepository.Authenticate")
	defer span.Finish()
	user, err := p.GetByEmail(ctx, email)
	if err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "userRepository.GetByEmail")
	}
	if err = user.Authenticate(password); err != nil {
		p.traceErr(span, err)
		return &entities.User{}, errors.Wrap(err, "user.Authenticate")
	}
	return user, nil
}

func (p *userRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
