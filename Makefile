GO ?= GO111MODULE=on GOOS=linux CGO_ENABLED=0 go
EXEC_NAME ?= afterburner-exporter
BUILDX ?= docker buildx
IMG_TOOLS ?= $(BUILDX) imagetools
PLATFORMS ?= linux/amd64,linux/i386,linux/arm64,linux/arm/v7
pkgs     = $(shell $(GO) list ./... | grep -v vendor)

DOCKER_IMAGE ?= kennedyoliveira/afterburner-exporter

ifeq ($(DOCKER_TAG),)
	DOCKER_TAG = latest
endif

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
	@echo ">>> compiling for windows x32"
	@GOOS=windows GOARCH=386 $(GO) build -v -o bin/$(EXEC_NAME)_win_x32.exe .

	@echo ">>> compiling for windows x64"
	@GOOS=windows GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_win_x64.exe .

	@echo ">>> compiling for linux x32"
	@GOOS=linux GOARCH=386 $(GO) build -v -o bin/$(EXEC_NAME)_unix_x32 .

	@echo ">>> compiling for linux x64"
	@GOOS=linux GOARCH=amd64 $(GO) build -v -o bin/$(EXEC_NAME)_unix_x64 .

	@echo ">>> compiling for arm"
	@GOOS=linux GOARCH=arm $(GO) build -v -o bin/$(EXEC_NAME)_unix_arm .

	@echo ">>> compiling for arm 64"
	@GOOS=linux GOARCH=arm64 $(GO) build -v -o bin/$(EXEC_NAME)_unix_arm64 .

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
	@echo ">> building multi-arch docker images, tag=$(DOCKER_TAG)"
	@$(BUILDX) build -f Dockerfile-cross \
			  --platform $(PLATFORMS) \
			  --tag $(DOCKER_IMAGE):$(DOCKER_TAG) \
			  --push \
			  .

docker-clean:
	@$(BUILDX) rm afterburner-exporter-builder

docker-inspect:
	@$(IMG_TOOLS) inspect $(DOCKER_IMAGE):$(DOCKER_TAG)
