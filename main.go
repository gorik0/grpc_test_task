package main

import (
	"github.com/joho/godotenv"
	"grpc/pkg/user/storage/postgres"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	connStr := os.Getenv("POSTGRES_URL")
	_, err = postgres.NewPostgres(connStr)
	if err != nil {
		panic(err)

	}
}
