package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	POSTGRES_URL  string
	GRPC_ADDR     string
	KAFKA_BROKERS []string
	TTL           time.Duration
	REDIS_ADDR    string
}

func NewConfig() *Config {
	ttlString := os.Getenv("TTL")
	ttl, _ := strconv.Atoi(ttlString)
	return &Config{
		POSTGRES_URL:  os.Getenv("POSTGRES_URL"),
		GRPC_ADDR:     os.Getenv("GRPC_ADDR"),
		KAFKA_BROKERS: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		TTL:           time.Duration(ttl),
		REDIS_ADDR:    os.Getenv("REDIS_ADDR"),
	}
}
