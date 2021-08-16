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

### Async + Task/Job Runner

Async process cukup membantu menghandle process yang panjang.
Dengan bantuan jaeger kita memliki visibility prosessnya.
Tetapi kita tidak memiliki control terhadap process yang berjalan.
Kita perlu tools untuk membantu melihat dan mengkrontrol setiap process.

Ada beberapa tools yang bisa digunakan. Misalnya [gocraft/work](https://github.com/gocraft/work).

Fitur yang ditawarkan:

* Fast and efficient. Faster than [this](https://www.github.com/jrallison/go-workers), [this](https://www.github.com/benmanns/goworker), and [this](https://www.github.com/albrow/jobs). See below for benchmarks.
* Reliable - don't lose jobs even if your process crashes.
* Middleware on jobs -- good for metrics instrumentation, logging, etc.
* If a job fails, it will be retried a specified number of times.
* Schedule jobs to happen in the future.
* Enqueue unique jobs so that only one job with a given name/arguments exists in the queue at once.
* Web UI to manage failed jobs and observe the system.
* Periodically enqueue jobs on a cron-like schedule.
* Pause / unpause jobs and control concurrency within and across processes

> Untuk mengunakan `gocraft/work` kita perlu service penerima dan processor. 
> Sehingga kita perlu [merefactor](#refactor) code yang kita buat.

### Refactor

Kita gunakan [cobra](github.com/spf13/cobra) untuk mempermudah membuat cli app.

```
cobra init --pkg-name github.com/zeihanaulia/go-long-process-api . -a "POC Long Process API"
cobra add api
cobra add worker
```

Menambahkan folder `presenters`, ambil term dari [the-clean-architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).
Didalam folder presenters, bikin folder `api` karena yang ingin ditampilkan adalah json api. Nanti hasilnya menjadi seperti ini:

```bash
- presenters
-- api
---- api.go
---- middleware.go
---- updater.go
```

- File `updater.go` akan diisi function `updaterHandler` yang sebelumnya ada pada `main.go`
- File `middleware.go` akan diisi function `traceMiddleware` yang sebelumnya ada pada `main.go`
- File `app.go` akan diisi dengan code yang ada dimain. dengan sedikin modifikasi.

Terakhir, Pada file `cmd/api.go` akan memanggil file app.go

```go
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Running API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		tracer, closer, err := tracing.Init("poc-task-processor")
		if err != nil {
			panic(fmt.Errorf("cannot start server %v", err))
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		return api.NewAPI(tracer).Run()
	},
}
```


## Simulate Process

### Run Application

```go
make api
make worker
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