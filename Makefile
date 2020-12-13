SRCS := $(shell find . -name '*.go')

.PHONY: all
all: test build

.PHONY: vet
vet:
	go vet ./...

.PHONY: golint
golint:
	for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done


.PHONY: lint
lint: golint vet

.PHONY: test
test: lint
	go test -race ./...

.PHONY: clean
clean:
	go clean -i ./...

build:
	go build -a -installsuffix cgo -o app .
