LOCAL_BIN:=${HOME}/.local/bin
PROG_NAME:=sc-loader
export GOEXPERIMENT=jsonv2

fix:
	@go fix ./...

build: fix
	@go build -ldflags "-s -w" .

test:
	@go test ./...

cover:
	@go test -race -coverprofile=coverage.out ./client/... ./manager/... ./models/... ./utils/...
	@go tool cover -func=coverage.out | tail -1

install: build
	$(info "installpath=${LOCAL_BIN}", progname=${PROG_NAME})
	@mv --verbose -t ${LOCAL_BIN} ${PROG_NAME}

uninstall:
	$(info "uninstall=${LOCAL_BIN}", progname=${PROG_NAME})
	@rm $(which ${PROG_NAME})
