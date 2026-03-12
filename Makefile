LOCAL_BIN:=${HOME}/.local/bin
PROG_NAME:=sc-loader


fix:
	@go fix ./...

build: fix
	@go build -ldflags "-s -w" .

install: build
	$(info "installpath=${LOCAL_BIN}", progname=${PROG_NAME})
	@mv --verbose -t ${LOCAL_BIN} ${PROG_NAME}

uninstall:
	$(info "uninstall=${LOCAL_BIN}", progname=${PROG_NAME})
	@rm $(which ${PROG_NAME})
