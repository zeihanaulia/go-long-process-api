package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

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
