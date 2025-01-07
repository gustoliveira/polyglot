build:
	@echo "Building and installing..."
	@go build -o att main.go
	@go install

run:
	@go run main.go

test:
	@echo "Testing..."
	@go test ./... -v

coverage:
	@echo "Calculating test coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

clean:
	@echo "Cleaning..."
	@rm -f main
	@rm -rf coverage.out

.PHONY: build run test clean
