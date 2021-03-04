sourcefiles = $(wildcard **/*.go)

build: $(sourcefiles)
	go build -o freecloud-server ./cmd/freecloud-server

run: build
	./freecloud-server

gocal:
	gocal

test: generatemock
	go test ./... --cover

generatemock:
	go install github.com/golang/mock/mockgen@v1.5.0
	go generate ./...