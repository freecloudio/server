sourcefiles = $(wildcard **/*.go)

build: $(sourcefiles)
	go build -o freecloud-server ./cmd/freecloud-server

run: build
	./freecloud-server

gocal:
	gocal

test:
	go test ./... --cover