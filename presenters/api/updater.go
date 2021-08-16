package api

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/zeihanaulia/go-long-process-api/pkg/response"
	"github.com/zeihanaulia/go-long-process-api/service"
)

func updaterHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span, _ := opentracing.StartSpanFromContext(ctx, "handler.product.updater")
	defer span.Finish()

	// Process
	go func(span opentracing.Span) {
		ctx := context.Background() // recreate context for avoid cancelation
		ctx = opentracing.ContextWithSpan(ctx, span)

		svc := service.NewService()
		if err := svc.Updater(ctx, &service.UpdaterRequest{}); err != nil {
			response.NewJSONResponse().SetError(response.ErrInternalServer)
			return
		}
	}(span)

	var traceID string
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		traceID = sc.TraceID().String()
	}

	resp := struct {
		Status  string `json:"status"`
		TraceID string `json:"trace_id"`
	}{
		Status:  "ok",
		TraceID: traceID,
	}
	response.NewJSONResponse().SetBody(resp).WriteResponse(rw)
}
