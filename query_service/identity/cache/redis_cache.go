package cache

import (
	"context"
	"encoding/json"
	"github.com/JECSand/identity-service/pkg/logging"
	"github.com/JECSand/identity-service/query_service/config"
	"github.com/JECSand/identity-service/query_service/identity/entities"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

const (
	redisUserMembershipPrefixKey  = "query:userMembership"
	redisGroupMembershipPrefixKey = "query:groupMembership"
	redisMembershipPrefixKey      = "query:membership"
	redisGroupPrefixKey           = "query:group"
	redisUserPrefixKey            = "query:user"
	redisTokenPrefixKey           = "query:token"
)

type redisCache struct {
	log         logging.Logger
	cfg         *config.Config
	redisClient redis.UniversalClient
}

func NewRedisCache(log logging.Logger, cfg *config.Config, redisClient redis.UniversalClient) *redisCache {
	return &redisCache{
		log:         log,
		cfg:         cfg,
		redisClient: redisClient,
	}
}

func (r *redisCache) PutUserMembership(ctx context.Context, key string, user *entities.UserMembership) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.PutUserMembership")
	defer span.Finish()
	b, err := json.Marshal(user)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}
	if err = r.redisClient.HSetNX(ctx, r.getRedisPrefixKey("userMembership"), key, b).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisPrefixKey("userMembership"), key)
}

func (r *redisCache) GetUserMembership(ctx context.Context, key string) (*entities.UserMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.GetUserMembership")
	defer span.Finish()
	b, err := r.redisClient.HGet(ctx, r.getRedisPrefixKey("userMembership"), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var user entities.UserMembership
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisPrefixKey("userMembership"), key)
	return &user, nil
}

func (r *redisCache) DeleteUserMembership(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteUserMembership")
	defer span.Finish()
	if err := r.redisClient.HDel(ctx, r.getRedisPrefixKey("userMembership"), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisPrefixKey("userMembership"), key)
}

func (r *redisCache) DeleteAllUserMemberships(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteAllUserMemberships")
	defer span.Finish()
	if err := r.redisClient.Del(ctx, r.getRedisPrefixKey("userMembership")).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisPrefixKey("userMembership"))
}

func (r *redisCache) PutGroupMembership(ctx context.Context, key string, user *entities.GroupMembership) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.PutGroupMembership")
	defer span.Finish()
	b, err := json.Marshal(user)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}
	if err = r.redisClient.HSetNX(ctx, r.getRedisPrefixKey("groupMembership"), key, b).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisPrefixKey("groupMembership"), key)
}

func (r *redisCache) GetGroupMembership(ctx context.Context, key string) (*entities.GroupMembership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.GetGroupMembership")
	defer span.Finish()
	b, err := r.redisClient.HGet(ctx, r.getRedisPrefixKey("groupMembership"), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var user entities.GroupMembership
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisPrefixKey("groupMembership"), key)
	return &user, nil
}

func (r *redisCache) DeleteGroupMembership(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteGroupMembership")
	defer span.Finish()
	if err := r.redisClient.HDel(ctx, r.getRedisPrefixKey("groupMembership"), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisPrefixKey("groupMembership"), key)
}

func (r *redisCache) DeleteAllGroupMemberships(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteAllMemberships")
	defer span.Finish()
	if err := r.redisClient.Del(ctx, r.getRedisPrefixKey("groupMembership")).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisPrefixKey("groupMembership"))
}

func (r *redisCache) PutMembership(ctx context.Context, key string, user *entities.Membership) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.PutMembership")
	defer span.Finish()
	b, err := json.Marshal(user)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}
	if err = r.redisClient.HSetNX(ctx, r.getRedisPrefixKey("membership"), key, b).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisPrefixKey("membership"), key)
}

func (r *redisCache) GetMembership(ctx context.Context, key string) (*entities.Membership, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.GetMembership")
	defer span.Finish()
	b, err := r.redisClient.HGet(ctx, r.getRedisPrefixKey("membership"), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var user entities.Membership
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisPrefixKey("membership"), key)
	return &user, nil
}

func (r *redisCache) DeleteMembership(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteMembership")
	defer span.Finish()
	if err := r.redisClient.HDel(ctx, r.getRedisPrefixKey("membership"), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisPrefixKey("membership"), key)
}

func (r *redisCache) DeleteAllMemberships(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteAllMemberships")
	defer span.Finish()
	if err := r.redisClient.Del(ctx, r.getRedisPrefixKey("membership")).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisPrefixKey("membership"))
}

func (r *redisCache) PutGroup(ctx context.Context, key string, user *entities.Group) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.PutGroup")
	defer span.Finish()
	b, err := json.Marshal(user)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}
	if err = r.redisClient.HSetNX(ctx, r.getRedisPrefixKey("group"), key, b).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisPrefixKey("group"), key)
}

func (r *redisCache) GetGroup(ctx context.Context, key string) (*entities.Group, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.GetGroup")
	defer span.Finish()
	b, err := r.redisClient.HGet(ctx, r.getRedisPrefixKey("group"), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var user entities.Group
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisPrefixKey("group"), key)
	return &user, nil
}

func (r *redisCache) DeleteGroup(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteGroup")
	defer span.Finish()
	if err := r.redisClient.HDel(ctx, r.getRedisPrefixKey("group"), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisPrefixKey("group"), key)
}

func (r *redisCache) DeleteAllGroups(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteAllGroups")
	defer span.Finish()
	if err := r.redisClient.Del(ctx, r.getRedisPrefixKey("group")).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisPrefixKey("group"))
}

func (r *redisCache) PutUser(ctx context.Context, key string, user *entities.User) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.PutUser")
	defer span.Finish()
	b, err := json.Marshal(user)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}
	if err = r.redisClient.HSetNX(ctx, r.getRedisPrefixKey("user"), key, b).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisPrefixKey("user"), key)
}

func (r *redisCache) GetUser(ctx context.Context, key string) (*entities.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.GetUser")
	defer span.Finish()
	b, err := r.redisClient.HGet(ctx, r.getRedisPrefixKey("user"), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var user entities.User
	if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisPrefixKey("user"), key)
	return &user, nil
}

func (r *redisCache) DeleteUser(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteUser")
	defer span.Finish()
	if err := r.redisClient.HDel(ctx, r.getRedisPrefixKey("user"), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisPrefixKey("user"), key)
}

func (r *redisCache) DeleteAllUsers(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteAllUsers")
	defer span.Finish()
	if err := r.redisClient.Del(ctx, r.getRedisPrefixKey("user")).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisPrefixKey("user"))
}

func (r *redisCache) PutToken(ctx context.Context, key string, blacklist *entities.Blacklist) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.PutToken")
	defer span.Finish()
	b, err := json.Marshal(blacklist)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}
	if err = r.redisClient.HSetNX(ctx, r.getRedisPrefixKey("token"), key, b).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisPrefixKey("token"), key)
}

func (r *redisCache) GetToken(ctx context.Context, key string) (*entities.Blacklist, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.GetToken")
	defer span.Finish()
	b, err := r.redisClient.HGet(ctx, r.getRedisPrefixKey("token"), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var bl entities.Blacklist
	if err = json.Unmarshal(b, &bl); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisPrefixKey("token"), key)
	return &bl, nil
}

func (r *redisCache) DeleteToken(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteToken")
	defer span.Finish()
	if err := r.redisClient.HDel(ctx, r.getRedisPrefixKey("token"), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisPrefixKey("token"), key)
}

func (r *redisCache) DeleteAllTokens(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "redisCache.DeleteAllTokens")
	defer span.Finish()
	if err := r.redisClient.Del(ctx, r.getRedisPrefixKey("token")).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisPrefixKey("token"))
}

func (r *redisCache) getRedisPrefixKey(prefixType string) string {
	switch prefixType {
	case "userMembership":
		if r.cfg.ServiceSettings.RedisUserMembershipPrefixKey != "" {
			return r.cfg.ServiceSettings.RedisUserMembershipPrefixKey
		}
		return redisUserMembershipPrefixKey
	case "groupMembership":
		if r.cfg.ServiceSettings.RedisGroupMembershipPrefixKey != "" {
			return r.cfg.ServiceSettings.RedisGroupMembershipPrefixKey
		}
		return redisGroupMembershipPrefixKey
	case "membership":
		if r.cfg.ServiceSettings.RedisMembershipPrefixKey != "" {
			return r.cfg.ServiceSettings.RedisMembershipPrefixKey
		}
		return redisMembershipPrefixKey
	case "group":
		if r.cfg.ServiceSettings.RedisGroupPrefixKey != "" {
			return r.cfg.ServiceSettings.RedisGroupPrefixKey
		}
		return redisGroupPrefixKey
	case "user":
		if r.cfg.ServiceSettings.RedisUserPrefixKey != "" {
			return r.cfg.ServiceSettings.RedisUserPrefixKey
		}
		return redisUserPrefixKey
	case "token":
		if r.cfg.ServiceSettings.RedisTokenPrefixKey != "" {
			return r.cfg.ServiceSettings.RedisTokenPrefixKey
		}
		return redisTokenPrefixKey
	}
	return "query"
}
