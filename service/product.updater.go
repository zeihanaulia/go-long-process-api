package service

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
)

func (s *service) Updater(ctx context.Context, payload *UpdaterRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.product.updater")
	defer span.Finish()

	// Update to DB
	if err := s.updateToDB(ctx); err != nil {
		return err
	}

	// Update to Elasticsearch
	if err := s.updateToElastic(ctx); err != nil {
		return err
	}

	return nil
}

func (s *service) updateToDB(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.product.updateToDB")
	defer span.Finish()

	// Simulate process db
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (s *service) updateToElastic(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.product.updateToElastic")
	defer span.Finish()

	// Simulate process elasticsearch
	time.Sleep(500 * time.Millisecond)
	return nil
}
