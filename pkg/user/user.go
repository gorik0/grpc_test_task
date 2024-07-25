package user

import (
	grpc2 "google.golang.org/grpc"
	"grpc/pkg/user/grpc"
	"grpc/pkg/user/grpc/userservice"
	"grpc/pkg/user/service"
)

func RegistrationUserServiceServer(registrator grpc2.ServiceRegistrar, services service.Services) *grpc.GRPC {
	g := &grpc.GRPC{
		Services: services,
	}
	userservice.RegisterUserServiceServer(registrator, g)
	return g
}
