GO111MODULE=auto

CONTAINER_NAME := jfreeland/trace-server:latest
DOCKER_OPTS := \
	-p 127.0.0.1:50000:8080/tcp \
	--rm --interactive --tty

.PHONY: build
build: ## Do stuff
	docker build . -t ${CONTAINER_NAME}

.PHONY: up
up: build ## Do stuff
	docker run ${DOCKER_OPTS} ${CONTAINER_NAME}

.PHONY: start
start: ## Do stuff
	curl -X POST http://localhost:50000/v0/start -H 'TraceHost: google.com'

.PHONY: graph
graph: ## Do stuff
	curl -X GET http://localhost:50000/v0/graph -H 'TraceHost: google.com' | jq .

.PHONY: stop
stop: ## Do stuff
	curl -X DELETE http://localhost:50000/v0/stop -H 'TraceHost: google.com'
