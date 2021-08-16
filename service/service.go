package service

import "context"

type Service interface {
	Updater(ctx context.Context, payload *UpdaterRequest) error
}

type service struct{}

func NewService() *service {
	return &service{}
}
