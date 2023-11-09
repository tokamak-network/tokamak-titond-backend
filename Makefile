BUILD_DIR=./build/bin

TARGET=$(BUILD_DIR)/titond

.PHONY: all run build titond image image-arm clean

all: run

run: $(TARGET)
	$(TARGET)

$(BUILD_DIR):
	mkdir -p $@

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
	docker build --build-arg TARGETARCH=amd64 --build-arg TARGETOS=linux -t titond-backend .

image-arm:
	docker build --build-arg TARGETARCH=arm64 --build-arg TARGETOS=linux -t titond-backend .


clean: $(BUILD_DIR)
	rm -rf $<
