SRC=$(shell find . -name '*.go')

plummet: $(SRC)
	go mod tidy
	go build ./cmd/plummet

release:
	we goreleaser --clean

.PHONY: goreleaser
goreleaser: nfpm
	go install github.com/goreleaser/goreleaser@latest

.PHONY: nfpm
nfpm:
	go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
