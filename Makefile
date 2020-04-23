SHELL=/bin/bash

callbackserver:
	docker build -f build/callbackserver/callbackserver.Dockerfile -t james65535/callbackserver:latest .
	docker push james65535/callbackserver:latest
	kubectl apply -f build/callbackserver/callbackserver.yaml

kafka:
	kubectl apply -f build/kafka/kafka.yaml

twitchclient:
	go run cmd/twitchclient/main.go

userconsumer:
	docker build -f build/userconsumer/userconsumer.Dockerfile -t james65535/userconsumer:latest .
	docker push james65535/userconsumer:latest
	kubectl apply -f build/userconsumer/userconsumer.yaml

all: kafka callbackserver userconsumer