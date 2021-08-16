# Long Process API

## Background Context

Service product updater yang dibuat tidak reliable.
Prosesnya pun semakin lambat karena penggunanya makin lama makin banyak.
Karena database dibagi dengan process lain, sehingga kerja database semakin tinggi dan membuat overload cpu > 98%.
Hal ini menyebabkan API yang diconsume oleh client menjadi lambat atau parahnya akan kena timeout dari client, padahal process masih berjalan.

## Approach Solutions

- [Async (Background Processing)](https://github.com/zeihanaulia/go-long-process-api/tree/01-async)
- Async + Task/Job Runner
- Async + Queue Messaging System

### Async (Background Processing)

Manfaatkan goroutine untuk membuat process berjalan dibackground.

```go
go func(span opentracing.Span) {
		ctx := context.Background() // recreate context for avoid cancelation
		ctx = opentracing.ContextWithSpan(ctx, span)

		svc := service.NewService()
		if err := svc.Updater(ctx, &service.UpdaterRequest{}); err != nil {
			response.NewJSONResponse().SetError(response.ErrInternalServer)
			return
		}
}(span)
```

Context dibuat ulang agar tidak kena cancelation dari prosess parentnya.
Inject span dari parent untuk keperluan tracing.
Expose tracing id untuk mencari lewat jaeger.

```go
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
```
## Simulate Process

### Run Application

```go
make run
make simulate
```

### Install Jaeger

Using Jaeger for visualize our request

```bash
docker run \
  --rm \
  --name jaeger \
  -p6831:6831/udp \
  -p16686:16686 \
  jaegertracing/all-in-one:latest
```

## Referensi

- https://en.wikipedia.org/wiki/Asynchrony_(computer_programming)