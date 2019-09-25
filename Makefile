NAME = chroma
ROOT = github.com/phR0ze/${NAME}
IMPORT = ${ROOT}/internal/${NAME}
VERSION := ${strip ${shell sed -En 's/version=(.*)/\1/p' VERSION}}
GIT_BRANCH := ${strip ${shell git rev-parse --abbrev-ref HEAD 2>/dev/null}}
GIT_COMMIT := ${strip ${shell git rev-parse --short HEAD 2>/dev/null}}
GIT_COMMIT_LONG := ${strip ${shell git rev-parse HEAD 2>/dev/null}}
GCFLAGS := -gcflags "all=-trimpath=${GOPATH}/src" -asmflags "all=-trimpath=${GOPATH}/src"
GOFLAGS := -ldflags '-X ${IMPORT}.VERSION=${VERSION} -X ${IMPORT}.GITCOMMIT=${GIT_COMMIT_LONG} -X ${IMPORT}.BUILDDATE=${shell date +%Y.%m.%d}'

.PHONY: build test

default: ${NAME}
${NAME}:
	@echo "Building..."
	@echo "------------------------------------------------------------------------"
	go build ${GOFLAGS} -o bin/${NAME} ./cmd/${NAME}
	
test:
	@echo "Testing..."
	@echo "------------------------------------------------------------------------"
	go test ./internal/${NAME}

clean:
	@rm -rf bin
	@mkdir bin
