export GOPATH := $(CURDIR)

guard:
	@echo "Build bChat server..."
	go build -o bin/bchat_guard bChat/guard

deps:
	@echo "Install Installing dependencies"
	@go get -u github.com/golang/dep/cmd/dep
	cd src/bChat; ${GOPATH}/bin/dep init; ${GOPATH}/bin/dep ensure -v