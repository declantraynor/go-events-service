APP_IMAGE = go-events-service
APP_CONTAINER = go-events-service-app
REDIS_CONTAINER = go-events-service-redis

RUN_REDIS = docker run --name $(REDIS_CONTAINER) -d redis

redis:
	if ! docker ps -a | grep $(REDIS_CONTAINER); then $(RUN_REDIS); fi
	if ! docker ps | grep $(REDIS_CONTAINER); then docker start $(REDIS_CONTAINER); fi

build:
	docker build -t $(APP_IMAGE) . 

run: build redis
	docker run -itP --link $(REDIS_CONTAINER):redis --rm --name $(APP_CONTAINER) $(APP_IMAGE)

test:
	go test -cover ./...
