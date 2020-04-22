GO ?= GO111MODULE=on go
EXEC_NAME ?= prometheus-msi-afterburner-exporter
DOCKER_IMAGE ?= kennedyoliveira/prometheus-msi-afterburner-exporter
BUILDX ?= docker buildx
IMG_TOOLS ?= $(BUILDX) imagetools
PLATFORMS ?= linux/amd64,linux/i386,linux/arm64,linux/arm/v7

.PHONY: build test compile clean docker-init docker-build docker-clean docker-inspect

all: test build

build:
	echo Building...
	$(GO) build -v -o bin/$(EXEC_NAME) .

test:
	echo Running tests...
	$(GO) test -v ./...

compile:
	GOOS=windows GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_winx64.exe .
	GOOS=linux GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_unix_x64 .
	GOOS=linux GOARCH=arm $(GO) build -v -o bin/$(EXEC_NAME)_unix_arm .

clean:
	rm -rf ./bin

docker-init:
	$(BUILDX) create --name afterburner-exporter-builder
	$(BUILDX) use afterburner-exporter-builder
	$(BUILDX) inspect --bootstrap afterburner-exporter-builder

docker-build: clean
	$(BUILDX) build -f Dockerfile-cross \
			  --platform $(PLATFORMS) \
			  --tag $(DOCKER_IMAGE) \
			  --push \
			  .

docker-clean:
	$(BUILDX) rm afterburner-exporter-builder

docker-inspect:
	$(IMG_TOOLS) inspect $(DOCKER_IMAGE)
