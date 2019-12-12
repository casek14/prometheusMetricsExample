APP_NAME := "prometheus-example-app"

.PHONY: help

help:
	@echo "help		Print help message."
	@echo "clean		Clean workspace."
	@echo "build-app	Build executable binary."
	@echo "build-container	Build docker container."

clean:
	rm -rf $(APP_NAME)

build-app: clean
	go build -o $(APP_NAME)

build-container:
	docker build -t quay.io/casek14/$(APP_NAME) -f Dockerfile .	
