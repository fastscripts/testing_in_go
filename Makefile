coverage:
	@echo "Running coverage tests..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"
	@rm coverage.out

test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests completed."	