NAME = chroma
ROOT = github.com/phR0ze/$(NAME)
IMPORT = github.com/phR0ze/$(NAME)/internal/$(NAME)
VERSION := $(strip $(shell sed -n 's/version=\(.*\)/\1/p' VERSION))
GIT_BRANCH := $(strip $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))
GIT_COMMIT := $(strip $(shell git rev-parse --short HEAD 2>/dev/null))
GIT_COMMIT_LONG := $(strip $(shell git rev-parse HEAD 2>/dev/null))
GCFLAGS := -gcflags "all=-trimpath=$(GOPATH)/src" -asmflags "all=-trimpath=$(GOPATH)/src"
GOFLAGS := -mod=vendor -ldflags '-X $(IMPORT).VERSION=$(VERSION) -X $(IMPORT).GITCOMMIT=$(GIT_COMMIT_LONG) -X $(IMPORT).BUILDDATE=$(shell date +%Y.%m.%d)'

.PHONY: build test

default: $(NAME)
$(NAME): vendor
	@echo "Building..."
	@echo "------------------------------------------------------------------------"
	go build ${GOFLAGS} -o bin/$(NAME) $(ROOT)/cmd/$(NAME)
	
vendor:
	go mod vendor

test:
	@echo "Testing..."
	@echo "------------------------------------------------------------------------"
	go test ${GOFLAGS} $(ROOT)/internal/$(NAME)

clean:
	@rm -rf bin
	@mkdir bin
