
help:
	@echo 'make build   - build the clipper executable'

collect: cmd/collect/main.go
	go build -o collect $^

build: collect
