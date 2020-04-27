run:
	go run ./cmd/bot

test:
	go test -count=1 -v -race ./...

DEFAULT_GOAL=run
