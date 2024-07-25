package grpc

import (
	"context"
	g "grpc/pkg/user/grpc/userservice"
	"grpc/pkg/user/model"
	"grpc/pkg/user/service"
	"log"
)

type GRPC struct {
	g.UnimplementedUserServiceServer
	Services service.Services
}

func (G *GRPC) CreateUser(ctx context.Context, u *g.User) (*g.Reply, error) {

	id, err := G.Services.CreateUser(ctx, u.Email)
	if err != nil {
		return &g.Reply{
			Id:     "",
			Status: err.Error(),
		}, err

	} else {
		log.Default().Println("CreateUser success")
		err = G.Services.LogNewUser(ctx)
		if err != nil {
			log.Default().Println("error while logging new user", err.Error())
		}
	}
	return &g.Reply{
		Id:     id,
		Status: "OK",
	}, nil
}

func (G *GRPC) DropUser(ctx context.Context, u *g.User) (*g.Reply, error) {
	err := G.Services.DeleteUser(ctx, u.Id)
	if err != nil {
		return &g.Reply{
			Status: err.Error(),
		}, err

	} else {
		return &g.Reply{
			Status: "DELETED",
		}, nil
	}

}
func (G *GRPC) GetUsers(ctx context.Context, u *g.SelectParams) (*g.UserList, error) {
	//	TRYin get user from cached then ( else from db and set to cache)
	var users []*model.User
	users, err := G.Services.GetUsersCached(ctx, int(u.Limit), int(u.Offset))
	switch {
	case err != nil:
		return nil, err
	case users == nil:
		{

			users, err = G.Services.GetUsers(ctx, int(u.Limit), int(u.Offset))
			if err != nil {
				return nil, err
			}
			err = G.Services.SetUsersCached(ctx, int(u.Limit), int(u.Offset), users)
			if err != nil {
				return nil, err
			}

		}

	}
	usersG := []*g.User{}

	for i, user := range users {
		usersG[i] = &g.User{
			Id:    user.ID,
			Email: user.Email,
		}

	}
	return &g.UserList{
		User: usersG,
	}, err
}
