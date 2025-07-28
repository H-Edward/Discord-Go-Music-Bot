BINARY_NAME=music-bot

# Choose the Go compiler
GOBUILD=go build
GO_SOURCE_HASH:=$(shell cat *.go | sha1sum | cut -c1-8)

all: build

build: 
	$(GOBUILD) -ldflags "-X 'main.GoSourceHash=$(GO_SOURCE_HASH)'" -o $(BINARY_NAME) -v

clean: 
	go clean
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

test: 
	go test -v ./...

docker-build:
	docker build -t $(BINARY_NAME) .

docker-network-create:
	docker network create musicbot-net || echo "Network already exists"

docker-run:
	docker run -d --name $(BINARY_NAME) --network musicbot-net --user 1000:1000 --read-only -v /app/config:/app/config:ro --cap-drop ALL \
	--security-opt no-new-privileges --memory=512M --cpus=3 --pids-limit=20 --restart unless-stopped $(BINARY_NAME)

docker-logs:
	docker logs -f $(BINARY_NAME)

docker-start:
	docker start $(BINARY_NAME)

docker-stop:
	docker stop $(BINARY_NAME)

docker-rm: 
	docker rm $(BINARY_NAME)

docker-rmi:
	docker rmi $(BINARY_NAME)

docker-kill:
	docker kill $(BINARY_NAME)



docker-network-rm:
	docker network rm musicbot-net || echo "Network does not exist"
