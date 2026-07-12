BINARY_NAME=music-bot

.PHONY: all build clean run test version up down logs restart
# Choose the Go compiler
GOBUILD=go build
GO_SOURCE_HASH:=$(shell find . -name "*.go" | sort | xargs cat | sha1sum | cut -c1-8)

all: build

build: 
	$(GOBUILD) -ldflags "-X 'discord-go-music-bot/internal/state.GoSourceHash=$(GO_SOURCE_HASH)'" -o $(BINARY_NAME) -v ./cmd/bot

run: build
	./$(BINARY_NAME)

test: 
	go test -v ./...

version:
	@echo "Version: $(GO_SOURCE_HASH)"

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

clean:
	docker-compose down --rmi all