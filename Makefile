APP_IMAGE = go-events-service
APP_CONTAINER = go-events-service-app

build:
	docker build -t $(APP_IMAGE) . 

run: build
	docker run -itP --rm --name $(APP_CONTAINER) $(APP_IMAGE)

test:
	go test -cover ./...
