package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/zeihanaulia/go-task-processor/pkg/response"
	"github.com/zeihanaulia/go-task-processor/pkg/tracing"
	"github.com/zeihanaulia/go-task-processor/service"
)

func main() {
	tracer, closer, err := tracing.Init("poc-task-processor")
	if err != nil {
		panic(fmt.Errorf("cannot start server %v", err))
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	http.Handle("/product", traceMiddleware(tracer, http.HandlerFunc(updaterHandler)))

	log.Println("listen on port http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func updaterHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler.product.updater")
	defer span.Finish()

	// Process
	svc := service.NewService()
	if err := svc.Updater(ctx, &service.UpdaterRequest{}); err != nil {
		response.NewJSONResponse().SetError(response.ErrInternalServer)
		return
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	response.NewJSONResponse().SetBody(resp).WriteResponse(rw)
}

func traceMiddleware(tracer opentracing.Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operationName := "HTTP " + r.Method + " " + r.URL.Path
		serverSpanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span, traceCtx := opentracing.StartSpanFromContextWithTracer(r.Context(), tracer, operationName, ext.RPCServerOption(serverSpanCtx))
		defer span.Finish()

		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.Path)

		// wraping untuk ambil status
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r.WithContext(traceCtx))

		status := ww.Status()
		ext.HTTPStatusCode.Set(span, uint16(status))

		if status >= 500 && status < 600 {
			ext.Error.Set(span, true)
			span.SetTag("error.type", fmt.Sprintf("%d: %s", status, http.StatusText(status)))
			span.LogKV(
				"event", "error",
				"message", fmt.Sprintf("%d: %s", status, http.StatusText(status)),
			)
		}
	})
}
