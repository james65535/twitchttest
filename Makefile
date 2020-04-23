SHELL=/bin/bash

twitchclient:
	go run cmd/twitchclient/main.go

callbackserver:
	docker build -f build/callbackserver/callbackserver.Dockerfile -t james65535/callbackserver:latest .
	docker push james65535/callback:latest
	kubectl apply -f build/callbackserver/callbackserver.yaml

all: callbackserver