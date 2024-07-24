package storage

import (
	"context"
	"grpc/pkg/user/model"
)

type CRUDL interface {
	CreateUser(ctx context.Context, user *model.UserDB) (int, error)

	GetAll(ctx context.Context, offset, limit int) ([]*model.UserDB, error)
	DeleteUser(ctx context.Context, id int) error
}
