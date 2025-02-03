build:
	@echo "Building and installing..."
	@go build -o polyglot main.go
	@go install
	@echo "Installing bash completion..."
	@polyglot completion bash > /tmp/polyglot-completion
	@bash -c 'source /tmp/polyglot-completion'

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
	@rm -rf dist/
	@rm -rf coverage.out

release:
	@echo "Building release with goreleaser locally..."
	@goreleaser check
	@goreleaser release --snapshot --clean

.PHONY: build run test covera clean release
