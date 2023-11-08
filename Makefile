.PHONY: all build titond

build:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	@echo "Done building"
	@echo "Run \"./build/bin/titond\" to launch titond backend."

titond:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	@echo "Done building"
	@echo "Run \"./build/bin/titond\" to launch titond backend."
