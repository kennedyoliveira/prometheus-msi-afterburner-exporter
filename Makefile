GO ?= GO111MODULE=on GOOS=linux CGO_ENABLED=1 go
EXEC_NAME ?= afterburner-exporter
DOCKER_IMAGE ?= kennedyoliveira/afterburner-exporter
BUILDX ?= docker buildx
IMG_TOOLS ?= $(BUILDX) imagetools
PLATFORMS ?= linux/amd64,linux/i386,linux/arm64,linux/arm/v7
pkgs     = $(shell $(GO) list ./... | grep -v vendor)

.PHONY: build test compile clean docker-init docker-build docker-clean docker-inspect format vet

all: test build

build:
	@echo ">> building"
	@$(GO) build -v -o bin/$(EXEC_NAME) .

test:
	@echo "Running tests"
	@$(GO) test -v $(pkgs)

compile:
	@echo ">> cross compiling"
	@GOOS=windows GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_winx64.exe .
	@GOOS=linux GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_unix_x64 .
	@GOOS=linux GOARCH=arm $(GO) build -v -o bin/$(EXEC_NAME)_unix_arm .

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

clean:
	@echo ">> cleaning"
	@rm -rf ./bin

docker-init:
	@$(BUILDX) create --name afterburner-exporter-builder
	@$(BUILDX) use afterburner-exporter-builder
	@$(BUILDX) inspect --bootstrap afterburner-exporter-builder

docker-build: clean
	@$(BUILDX) build -f Dockerfile-cross \
			  --platform $(PLATFORMS) \
			  --tag $(DOCKER_IMAGE) \
			  --push \
			  .

docker-clean:
	@$(BUILDX) rm afterburner-exporter-builder

docker-inspect:
	@$(IMG_TOOLS) inspect $(DOCKER_IMAGE)
