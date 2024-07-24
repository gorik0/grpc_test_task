package cache

import (
	"context"
	"grpc/pkg/user/model"
)

type Cacher interface {
	GetList(ctx context.Context, offset, limit int) ([]*model.User, error)
	SetList(ctx context.Context, users []*model.User, limit, offset int) error
}
