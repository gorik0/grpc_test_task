package service

import (
	"context"
	"grpc/pkg/user/cache"
	"grpc/pkg/user/logger"
	"grpc/pkg/user/model"
	"grpc/pkg/user/storage"
	"strconv"
)

type Services interface {
	CreateUser(ctx context.Context, email string) (id string, err error)
	DeleteUser(ctx context.Context, idUser string) error
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	GetUsersCached(ctx context.Context, limit, offset int) ([]*model.User, error)
	SetUsersCached(ctx context.Context, limit, offset int, users []*model.User) error
	LogNewUser(ctx context.Context) error
}
type Service struct {
	storage storage.CRUDL
	cacher  cache.Cacher
	logger  logger.Logger
}

func (s Service) CreateUser(ctx context.Context, email string) (id string, err error) {
	user := model.UserDB{
		Email: email,
	}

	idInt, err := s.storage.CreateUser(ctx, &user)
	if err != nil {
		return "", err

	}
	id = strconv.Itoa(idInt)
	return id, nil

}

func (s Service) DeleteUser(ctx context.Context, idUser string) error {
	id, _ := strconv.Atoi(idUser)
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	all, err := s.storage.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err

	}
	users := make([]*model.User, len(all))
	for i, user := range all {
		users[i] = &model.User{
			ID:    user.ID,
			Email: user.Email,
		}

	}
	return users, err
}

func (s Service) GetUsersCached(ctx context.Context, limit, offset int) ([]*model.User, error) {
	return s.cacher.GetList(ctx, limit, offset)

}

func (s Service) SetUsersCached(ctx context.Context, limit, offset int, users []*model.User) error {
	return s.cacher.SetList(ctx, users, limit, offset)
}

func (s Service) LogNewUser(ctx context.Context) error {
	return s.logger.LoNewUser(ctx)
}

func NewServices(storage storage.CRUDL, cacher cache.Cacher, logger logger.Logger) Services {
	return &Service{
		storage: storage,
		cacher:  cacher,
		logger:  logger,
	}
}
