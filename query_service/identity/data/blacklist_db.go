package data

import (
	"context"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/pkg/utilities"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// blacklistEntity structures a user BSON document to save in a users collection
type blacklistEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken string             `bson:"access_token,omitempty" validate:"required,min=3,max=500"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

// getID returns the unique identifier of the blacklistEntity
func (u *blacklistEntity) getID() primitive.ObjectID {
	return u.ID
}

// newBlacklistEntity initializes a new pointer to a blacklistEntity struct from a pointer to a JSON Blacklist struct
func newBlacklistEntity(u *entities.Blacklist) (um *blacklistEntity, err error) {
	um = &blacklistEntity{
		AccessToken: u.AccessToken,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if utilities.CheckID(u.ID) == nil {
		um.ID, err = utilities.LoadObjectIDString(u.ID)
	}
	return
}

// toRoot creates and return a new pointer to a models.Blacklist JSON struct from a pointer to a BSON userModel
func (u *blacklistEntity) toRoot() *entities.Blacklist {
	um := &entities.Blacklist{
		AccessToken: u.AccessToken,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	if utilities.CheckID(u.ID.Hex()) == nil {
		um.ID = utilities.LoadUUIDString(u.ID)
	}
	return um
}

func loadBlacklistEntities(ms []*blacklistEntity) (users []*entities.Blacklist) {
	for _, m := range ms {
		users = append(users, m.toRoot())
	}
	return
}

type blacklistRepository struct {
	log logging.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewBlacklistRepository(log logging.Logger, cfg *config.Config, db *mongo.Client) *blacklistRepository {
	return &blacklistRepository{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

func (p *blacklistRepository) BlacklistToken(ctx context.Context, bList *entities.Blacklist) (*entities.Blacklist, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistRepository.BlacklistToken")
	defer span.Finish()
	p.log.Info(bList.ID)
	ent, err := newBlacklistEntity(bList)
	if err != nil {
		p.traceErr(span, err)
		return &entities.Blacklist{}, errors.Wrap(err, "newBlacklistEntity")
	}
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Blacklist)
	_, err = collection.InsertOne(ctx, ent, &options.InsertOneOptions{})
	if err != nil {
		p.traceErr(span, err)
		return &entities.Blacklist{}, errors.Wrap(err, "InsertOne")
	}
	return bList, nil
}

func (p *blacklistRepository) CheckBlacklist(ctx context.Context, accessToken string) (*entities.Blacklist, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "blacklistRepository.CheckBlacklist")
	defer span.Finish()
	collection := p.db.Database(p.cfg.Mongo.DB).Collection(p.cfg.MongoCollections.Blacklist)
	var ent blacklistEntity
	if err := collection.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(&ent); err != nil {
		p.traceErr(span, err)
		return &entities.Blacklist{}, errors.Wrap(err, "Decode")
	}
	return ent.toRoot(), nil
}

func (p *blacklistRepository) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
