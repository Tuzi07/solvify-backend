test:
	go test ./... -v -coverprofile=coverage.cov -coverpkg=./...

run:
	go run ./cmd/main.go

clean:
	rm coverage.cov

build: test
	go build -o solvify-backend ./cmd/main.go
	make clean