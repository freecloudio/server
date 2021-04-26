sourcefiles = $(wildcard **/*.go)

build: $(sourcefiles)
	go build -o freecloud-server ./cmd/freecloud-server

run: build
	./freecloud-server

gocal:
	go install github.com/mheidinger/gocal
	gocal

test: generate
	go test ./... --cover

generate:
	go install github.com/golang/mock/mockgen
	go get github.com/99designs/gqlgen
	go generate ./...
