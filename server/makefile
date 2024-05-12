all: build test coverage-report

tidy:
	go mod tidy

build:
	go build -o bin/server ./main.go

assembly:
	go tool compile -S x.go

test:
	go test ./... --cover -coverprofile reports/coverage

coverage-report: test
	go tool cover -html reports/coverage