.PHONY: build run test vet clean

build:
	go build -o bro .

run:
	go run .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f bro
