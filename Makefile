TARGET=./build/bin/titond

.PHONY: all run build titond image image-arm clean

all: run

run: $(TARGET)
	$(TARGET)

$(TARGET): build

build:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	@echo "Done building"
	@echo "Run \"./build/bin/titond\" to launch titond backend."

titond:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/bin/titond ./cmd/titond/main.go
	@echo "Done building"
	@echo "Run \"./build/bin/titond\" to launch titond backend."

image:
	docker build --build-arg TARGETOS=linux -t titond-backend .

image-amd:
	docker build --build-arg TARGETARCH=amd64 --build-arg TARGETOS=linux -t titond-backend .

image-arm:
	docker build --build-arg TARGETARCH=arm64 --build-arg TARGETOS=linux -t titond-backend .


clean: 
	rm -rf $(TARGET)
