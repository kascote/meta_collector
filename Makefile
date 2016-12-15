MAJOR_VERSION := 0
MINOR_VERSION := 1
VERSION := ${MAJOR_VERSION}.${MINOR_VERSION}
#RELEASE := $(shell git describe --tags)

help:
	@echo 'make build   - build the collector executable'
	@echo 'make tag     - tag the current HEAD with VERSION'
	@echo 'make archive - create an archive of the current HEAD with VERSION'
	@echo 'make all     - tag, build and archive VERSION'

collect: cmd/collect/main.go
	go build -ldflags="-X main.version=${VERSION}" -o $@ $^

tag:
	git tag -a ${VERSION} -m "${VERSION} release"

archive: collect-${VERSION}.zip

collect-${VERSION}.zip: collect
	git archive -o $@ HEAD
	zip $@ collect

build: collect

all: tag build archive

.PHONY: clean
clean:
	rm -f collect collect-*.zip