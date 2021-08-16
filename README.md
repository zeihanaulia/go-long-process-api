# Long Process API

## Background Context

Service product updater yang dibuat tidak reliable.
Prosesnya pun semakin lambat karena penggunanya makin lama makin banyak.
Karena database dibagi dengan process lain, sehingga kerja database semakin tinggi dan membuat overload cpu > 98%.
Hal ini menyebabkan API yang diconsume oleh client menjadi lambat atau parahnya akan kena timeout dari client, padahal process masih berjalan.

## Approach Solutions

- Async (Background Processing)
- Async + Task/Job Runner
- Async + Queue Messaging System

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