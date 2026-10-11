LOCAL_BIN:=${HOME}/.local/bin
PROG_NAME:=sc-loader
GOLANGCI_LINT?=golangci-lint-v2

# golangci-lint fmt --stdin -E goimports,golines
fmt:
	@$(GOLANGCI_LINT) fmt -E goimports,golines ./...

fix: fmt
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
