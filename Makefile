.PHONY: build test run clean

build:
	go build -o savepoint main.go

test:
	go test ./...

run:
	go run main.go

clean:
	rm -f savepoint