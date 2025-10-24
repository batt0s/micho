build:
	go build -o ./bin/micho ./cmd/server

run:
	./bin/micho

dev:
	find . | grep '\.go$\' | entr -r go run ./cmd/server

