package queries

import (
	"context"
	"errors"
	"github.com/JECSand/identity-service/pkg/enums"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/cache"
	"github.com/JECSand/identity-service/query_service/identity/data"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/opentracing/opentracing-go"
)

// AuthenticateHandler ...
type AuthenticateHandler interface {
	Handle(ctx context.Context, query *AuthenticateQuery) (*entities.User, error)
}

type authenticateHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewAuthenticateHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *authenticateHandler {
	return &authenticateHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (q *authenticateHandler) Handle(ctx context.Context, query *AuthenticateQuery) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authenticateHandler.Handle")
	defer span.Finish()
	user, err := q.mongoDB.AuthenticateUser(ctx, query.Email, query.Password)
	if err != nil {
		q.log.WarnMsg("mongoDB.AuthenticateUser", err)
		return nil, err
	}
	return user, nil
}

// ValidateHandler ...
type ValidateHandler interface {
	Handle(ctx context.Context, query *ValidateQuery) (*entities.User, error)
}

type validateHandler struct {
	log        logging.Logger
	cfg        *config.Config
	mongoDB    data.Database
	redisCache cache.Cache
}

func NewValidateHandler(log logging.Logger, cfg *config.Config, mongoDB data.Database, redisCache cache.Cache) *validateHandler {
	return &validateHandler{
		log:        log,
		cfg:        cfg,
		mongoDB:    mongoDB,
		redisCache: redisCache,
	}
}

func (s *validateHandler) validateToken(ctx context.Context, query *ValidateQuery) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "validateHandler.validateToken")
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)
	userChan := make(chan *entities.User)
	defer func() {
		span.Finish()
		cancel()
		close(errChan)
		close(userChan)
	}()
	go func() { // Ensure token being validated is not in blacklist
		_, err := s.mongoDB.CheckTokenBlacklist(ctx, query.AccessToken)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err == nil {
			err = errors.New("token is blacklisted")
			s.log.WarnMsg("mongoDB.CheckTokenBlacklist", err)
		} else if err.Error() == "Decode: mongo: no documents in result" {
			err = nil
		}
		errChan <- err
	}()
	go func() { // Ensure userId references a valid user
		user, err := s.mongoDB.GetUserById(ctx, query.UserID)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err != nil {
			s.log.WarnMsg("mongoDB.GetUserById", err)
		}
		errChan <- err
		userChan <- user
	}()
	var errs []error
	var user *entities.User
	for i := 0; i < 3; i++ {
		select {
		case err := <-errChan:
			if err != nil {
				errs = append(errs, err)
			}
		case user = <-userChan:
		}
	}
	if len(errs) > 0 {
		return &entities.User{}, errs[0]
	}
	return user, nil
}

func (s *validateHandler) validatePassword(ctx context.Context, query *ValidateQuery) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "validateHandler.validatePassword")
	defer span.Finish()
	user, err := s.mongoDB.GetUserById(ctx, query.UserID)
	if err != nil {
		s.log.WarnMsg("mongoDB.GetUserById", err)
		return &entities.User{}, err
	}
	if err = user.Authenticate(query.AccessToken); err != nil {
		s.log.WarnMsg("mongoDB.GetUserById", err)
		return &entities.User{}, err
	}
	return user, nil
}

func (s *validateHandler) Handle(ctx context.Context, query *ValidateQuery) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "validateHandler.Handle")
	defer span.Finish()
	switch query.ValidationType {
	case enums.TOKEN:
		return s.validateToken(ctx, query)
	case enums.PASSWORD:
		return s.validatePassword(ctx, query)
	default:
		return &entities.User{}, errors.New("invalid ValidationType")
	}
}
