SRC=$(shell find . -name '*.go')

plummet: $(SRC)
	go mod tidy
	go build ./cmd/plummet
