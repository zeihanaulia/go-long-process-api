package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/zeihanaulia/go-task-processor/pkg/tracing"
	"github.com/zeihanaulia/go-task-processor/service"
	"syreclabs.com/go/faker"
)

func main() {
	tracer, closer, err := tracing.Init("simulate-task-processing")
	if err != nil {
		panic(fmt.Errorf("cannot start server %v", err))
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	request()
}

func request() {
	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "simulate.request")
	defer span.Finish()

	URL := "http://localhost:3000/product"
	var jsonStr = createPayload()
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	if span := opentracing.SpanFromContext(ctx); span != nil {
		_ = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func createPayload() []byte {
	var product = make([]service.Product, 0)
	for i := 0; i < 10; i++ {
		product = append(product, service.Product{
			ID:      faker.Number().NumberInt(11),
			StoreID: faker.Code().Abn(),
			SKU:     faker.Code().Rut(),
			Name:    faker.Commerce().ProductName(),
		})
	}

	payload := service.UpdaterRequest{Data: product}

	b, _ := json.Marshal(payload)
	return b
}
