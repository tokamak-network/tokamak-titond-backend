TARGET=./build/bin/titond

.PHONY: all run build titond image image-amd image-arm clean test check

all: run

run: $(TARGET)
	$<

$(TARGET): build

check:
	swag init -g cmd/titond/main.go -o api
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	./build/bin/titond check-swagger

build:
	swag init -g cmd/titond/main.go -o api
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	@echo "Done building"
	@echo "Run \"./build/bin/titond\" to launch titond backend."

titond:
	swag init -g cmd/titond/main.go -o api
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	@echo "Done building"
	@echo "Run \"./build/bin/titond\" to launch titond backend."

image:
	docker build --build-arg TARGETOS=linux -t titond-backend .

image-amd:
	docker build --build-arg TARGETARCH=amd64 --build-arg TARGETOS=linux -t titond-backend .

image-arm:
	docker build --build-arg TARGETARCH=arm64 --build-arg TARGETOS=linux -t titond-backend .

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	./scripts/check_coverage_unittest.sh

clean: 
	rm -rf $(TARGET)
	rm -rf docs
