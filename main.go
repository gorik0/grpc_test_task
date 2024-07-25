package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc/config"
	"grpc/pkg/user"
	"grpc/pkg/user/cache/redis"
	"grpc/pkg/user/logger/kafka"
	"grpc/pkg/user/service"
	"grpc/pkg/user/storage/postgres"
	"log"
	"net"
)

func main() {
	// :::
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg := config.NewConfig()

	crudl, err := postgres.NewPostgres(cfg.POSTGRES_URL)
	if err != nil {
		log.Fatalf("Error while NewPostgres  ::: ", err.Error())
	}
	looger, err := kafka.NewKafka(cfg.KAFKA_BROKERS)
	if err != nil {
		log.Fatalf("Error while NewKafka  ::: ", err.Error())
	}
	cacher, err := redis.NewRedisClient(cfg.REDIS_ADDR, cfg.TTL)
	if err != nil {
		log.Fatalf("Error while NewRedisClient  ::: ", err.Error())
	}

	services := service.NewServices(crudl, cacher, looger)

	grpcServer := grpc.NewServer()
	user.RegistrationUserServiceServer(grpcServer, services)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPC_ADDR)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
