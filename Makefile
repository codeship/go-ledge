.PHONY: \
	all \
	testall \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	generate \
	build \
	install \
	cov \
	test \
	clean

all: test install

testall: cov lint test

deps:
	go get -d -v ./...

updatedeps:
	go get -d -v -u -f ./...

testdeps:
	go get -d -v -t ./...
	go get -v golang.org/x/tools/cmd/vet
	go get -v github.com/kisielk/errcheck

updatetestdeps:
	go get -d -v -t -u -f ./...
	go get -v -u -f golang.org/x/tools/cmd/vet
	go get -v -u -f github.com/kisielk/errcheck

generate:
	go generate ./...

build: deps generate
	go build ./...

install: deps generate
	go install ./...

cov: testdeps generate
	go get -v github.com/axw/gocov/gocov
	go get golang.org/x/tools/cmd/cover
	gocov test | gocov report

lint: testdeps generate
	go get -v github.com/golang/lint/golint
	golint ./...

test: testdeps generate
	go test -test.v ./...
	go vet ./...
	errcheck ./...

clean:
	go clean -i ./...
