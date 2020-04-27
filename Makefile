SHELL=/bin/bash

callbackserver:
	docker build -f build/callbackserver/callbackserver.Dockerfile -t james65535/callbackserver:latest .
	docker push james65535/callbackserver:latest
	kubectl apply -f build/callbackserver/callbackserver.yaml

dispatcher:
	docker build -f build/dispatch/dispatcher.Dockerfile -t james65535/dispatcher:latest .
	docker push james65535/dispatcher:latest
	kubectl apply -f build/dispatcher/dispatcher.yaml

kafka:
	kubectl apply -f build/kafka/kafka.yaml

test:
	kind create cluster
	kubectl apply -f build/kafka/test-kafka.yaml
	kubectl port-forward service/kafka-svc 9092

twitchclient:
	go run cmd/twitchclient/main.go

userconsumer:
	docker build -f build/userconsumer/userconsumer.Dockerfile -t james65535/userconsumer:latest .
	docker push james65535/userconsumer:latest
	kubectl apply -f build/userconsumer/userconsumer.yaml

all: kafka callbackserver dispatcher userconsumer