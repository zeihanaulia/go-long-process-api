PROJECTNAME := $(shell basename "$(PWD)")
STIME 		:= $(shell date +%s)

.PHONY: build
build:
	@echo ">  Building Program..."
	go build -ldflags="-s -w" -o bin/${PROJECTNAME} main.go; 
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"

## run: start without docker
.PHONY: run
run: build
	@echo "  >  Starting Program..."
	./bin/${PROJECTNAME} api
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"

## simulate: start without docker
.PHONY: simulate
simulate:
	@echo "  >  Simulating Request..."
	go run ./fixtures/simulate.go
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"