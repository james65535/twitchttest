SHELL=/bin/bash

twitchclient:
	go run cmd/twitchclient/main.go
server:
	docker build -f build/Docker/server.Dockerfile -t james65535/twitchtest:latest .
	docker push james65535/twitchtest:latest
	kubectl apply -f build/server/server.yaml

all: server