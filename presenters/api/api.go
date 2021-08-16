package api

import (
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

type apis struct {
	tracer opentracing.Tracer
}

func NewAPI(tracer opentracing.Tracer) *apis {
	return &apis{tracer}
}

func (a *apis) Run() error {
	http.Handle("/product", traceMiddleware(a.tracer, http.HandlerFunc(updaterHandler)))
	log.Println("listen on port http://localhost:3000")
	return http.ListenAndServe(":3000", nil)
}
