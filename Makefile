GO ?= GO111MODULE=on go
EXEC_NAME ?= prometheus-msi-afterburner-exporter
DOCKER_IMAGE ?= kennedyoliveira/prometheus-msi-afterburner-exporter
BUILDX ?= docker buildx
IMG_TOOLS ?= $(BUILDX) imagetools

all: test build

.PHONY: build
build:
	echo Building...
	$(GO) build -v -o bin/$(EXEC_NAME) .

.PHONY: test
test:
	echo Running tests...
	$(GO) test -v ./...

.PHONY: compile
compile:
	GOOS=windows GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_winx64.exe .
	GOOS=linux GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_unix_x64 .
	GOOS=linux GOARCH=arm $(GO) build -v -o bin/$(EXEC_NAME)_unix_arm .

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: docker-prepare
docker-prepare:
	$(BUILDX) create --name afterburner-exporter-builder
	$(BUILDX) use afterburner-exporter-builder
	$(BUILDX) inspect --bootstrap afterburner-exporter-builder

.PHONY: docker-build
docker-build:
	$(BUILDX) build -f Dockerfile-cross --platform linux/386,linux/amd64,linux/arm/v7,linux/arm/v6,linux/arm64 --tag $(DOCKER_IMAGE) --push .

.PHONY: docker-cleanup
docker-cleanup:
	$(BUILDX) rm afterburner-exporter-builder

.PHONY: docker-inspect
docker-inspect:
	$(IMG_TOOLS) inspect $(DOCKER_IMAGE)
