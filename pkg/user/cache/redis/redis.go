package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"grpc/pkg/user/cache"
	"grpc/pkg/user/model"
	"time"
)

type RedisClient struct {
	redis *redis.Client
	ttl   time.Duration
}

func (r RedisClient) GetList(ctx context.Context, offset, limit int) ([]*model.User, error) {
	bytes, err := r.redis.Get(ctx, fmt.Sprintf("userList_l%d_o%d", limit, offset)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var users []*model.User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r RedisClient) SetList(ctx context.Context, users []*model.User, limit, offset int) error {
	err := r.redis.Set(ctx, fmt.Sprintf("userList_l%d_o%d", len(users), 0), users, r.ttl).Err()
	if err != nil {
		return err

	}
	return nil
}

func NewRedisClient(addr string, ttl time.Duration) (*RedisClient, error) {

	opts, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err

	}
	client := redis.NewClient(opts)
	return &RedisClient{
		redis: client,
		ttl:   ttl,
	}, nil

}

var _ cache.Cacher = &RedisClient{}
