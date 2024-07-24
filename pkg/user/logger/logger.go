package logger

import "context"

type Logger interface {
	LoNewUser(ctx context.Context) error
}

type Log struct {
	Logger
}
