GO ?= GO111MODULE=on CGO_ENABLED=0 go
EXEC_NAME ?= afterburner-exporter
BUILDX ?= docker buildx
IMG_TOOLS ?= $(BUILDX) imagetools
PLATFORMS ?= linux/amd64,linux/i386,linux/arm64,linux/arm/v7
pkgs     = $(shell $(GO) list ./... | grep -v vendor)
VER_SHA1 = $(shell git rev-parse HEAD)
BRANCH      = $(shell git rev-parse --abbrev-ref HEAD)
DOCKER_IMAGE ?= kennedyoliveira/afterburner-exporter
CUR_DATE = $(shell date +'%Y-%m-%d_%T')
GIT_TAG = $(shell git describe --tags)

ifeq ($(GIT_TAG),)
	GIT_TAG = DEV
endif

ifeq ($(DOCKER_TAG),)
	ifeq ($(BRANCH),develop)
	DOCKER_TAG = dev
	else
	DOCKER_TAG = latest
	endif
endif

.PHONY: build test compile clean docker-init docker-build docker-clean docker-inspect format vet info

all: test build

info:
	@echo "--- Go stuff ---"
	@echo "GO               -> $(GO)"
	@echo "pkgs             -> $(pkgs)"
	@echo "EXEC_NAME        -> $(EXEC_NAME)"

	@echo "--- Docker ---"
	@echo "DOCKER_IMAGE     -> $(DOCKER_IMAGE)"
	@echo "DOCKER_TAG       -> $(DOCKER_TAG)"
	@echo "PLATFORMS        -> $(PLATFORMS)"
	@echo "BUILDX           -> $(BUILDX)"
	@echo "IMG_TOOLS        -> $(IMG_TOOLS)"

	@echo "--- Git stuff ---"
	@echo "CUR_DATE         -> $(CUR_DATE)"
	@echo "GIT_TAG          -> $(GIT_TAG)"
	@echo "BRANCH           -> $(BRANCH)"
	@echo "VER_SHA1         -> $(VER_SHA1)"


build:
	@echo ">> building Branch=$(BRANCH), SHA1=$(VER_SHA1), TAG=$(GIT_TAG), Date=$(CUR_DATE)"
	@$(GO) build -v \
	  -o bin/$(EXEC_NAME) \
	  -ldflags "-X main.branch=$(BRANCH) -X main.sha1=$(VER_SHA1) -X main.buildDate=$(CUR_DATE) -X main.tag=$(GIT_TAG)" \
	  .

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
