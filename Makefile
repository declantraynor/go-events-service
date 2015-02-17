APP_IMAGE = go-events-service
TEST_IMAGE = go-events-service-test

APP_CONTAINER = go-events-service-app
REDIS_CONTAINER = go-events-service-redis
TEST_CONTAINER = go-events-service-test-container

RUN_REDIS = docker run --name $(REDIS_CONTAINER) -d redis

redis:
	# Builds a new redis container or starts an existing one
	if ! docker ps -a | grep $(REDIS_CONTAINER); then $(RUN_REDIS); fi
	if ! docker ps | grep $(REDIS_CONTAINER); then docker start $(REDIS_CONTAINER); fi

build:
	docker build -t $(APP_IMAGE) -f Dockerfile.build .

run: build redis
	docker run -itP --link $(REDIS_CONTAINER):redis --rm --name $(APP_CONTAINER) $(APP_IMAGE)

test:
	docker build -t $(TEST_IMAGE) -f Dockerfile.test .
	docker run -it --rm --name $(TEST_CONTAINER) $(TEST_IMAGE) go test -cover ./...
