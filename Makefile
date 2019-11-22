.PHONY: help check build clean

help:
	@echo "Targets:"
	@echo "- build: Compile main.go"
	@echo "- check: Format main.go"
	@echo "- test: Test main.go"
	@echo "- clean: Cleanup executables"

check:
	@go fmt main.go

build:
	@go build -o main main.go

test:
	@go test main_test.go main.go

clean:
	@rm -f main