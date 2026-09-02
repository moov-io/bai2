# generated-from:8ef700b33a05ab58ec9e7fd3ad1a0d8a99a742beeefc09d26bc7e4b6dd2ad699 DO NOT REMOVE, DO UPDATE

PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
PWD := $(shell pwd)

ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

ifndef VERSION
    VERSION := v.0.0.0
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

all: install update build

.PHONY: install
install:
	go mod tidy
	go mod vendor

update:
	go mod vendor

build:
	go build -mod=vendor -ldflags "-X github.com/moov-io/bai2.Version=${VERSION}" -o bin/bai2 github.com/moov-io/bai2/cmd/bai2

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=45.0 GOLANGCI_LINTERS=gosec ./lint-project.sh
endif

.PHONY: teardown
teardown:
	-docker-compose down --remove-orphans

docker: update docker-hub

docker-hub:
	docker build --pull --build-arg VERSION=${VERSION} -t moov/bai2:${VERSION} -f Dockerfile .
	docker tag moov/bai2:${VERSION} moov/bai2:latest

docker-push:
	docker push moov/bai2:${VERSION}
	docker push moov/bai2:latest

.PHONY: dev-docker
dev-docker: update
	docker build --pull --build-arg VERSION=${DEV_VERSION} -t moov/bai2:${DEV_VERSION} -f Dockerfile .

.PHONY: dev-push
dev-push:
	docker push moov/bai2:${DEV_VERSION}

# Extra utilities not needed for building

run: update build
	./bin/bai2

docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moov-io/bai2:${VERSION}

test: update
	go test -cover github.com/moov-io/bai2/...

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf cover.out coverage.txt misspell* staticcheck*
	@rm -rf ./bin/
endif

# For open source projects

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@

# Map uname -m to Go GOARCH (linux/windows dist). Darwin is always arm64.
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
  HOST_ARCH := amd64
else ifeq ($(UNAME_M),amd64)
  HOST_ARCH := amd64
else ifeq ($(UNAME_M),aarch64)
  HOST_ARCH := arm64
else ifeq ($(UNAME_M),arm64)
  HOST_ARCH := arm64
else
  HOST_ARCH := $(UNAME_M)
endif

# Tagged GitHub releases are built by GoReleaser (.goreleaser.yaml).
# This target is a local one-off for the host OS/arch.
dist: clean
ifeq ($(OS),Windows_NT)
	CGO_ENABLED=0 GOOS=windows go build -ldflags "-X github.com/moov-io/bai2.Version=${VERSION}" -o bin/bai2.exe github.com/moov-io/bai2/cmd/bai2
else ifeq ($(PLATFORM),darwin)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/moov-io/bai2.Version=${VERSION}" -o bin/bai2-darwin-arm64 github.com/moov-io/bai2/cmd/bai2
else
	CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(HOST_ARCH) go build -ldflags "-X github.com/moov-io/bai2.Version=${VERSION}" -o bin/bai2-$(PLATFORM)-$(HOST_ARCH) github.com/moov-io/bai2/cmd/bai2
endif
