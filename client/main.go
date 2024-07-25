package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/pkg/user/grpc/userservice"
)

func main() {
	dial, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("DIAL" + err.Error())
	}
	client := userservice.NewUserServiceClient(dial)
	userID, err := client.CreateUser(context.Background(), &userservice.User{
		Email: "d4",
	})
	if err != nil {
		panic("CREATE USER" + err.Error())

	}
	println("DONE!!! -->>> ", userID)
}
