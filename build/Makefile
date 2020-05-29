IMAGE_NAME := requester
CONTAINER_NAME := requester
URL := http://3.23.114.100:9898/newTask

build:
	cd .. && docker build -t $(IMAGE_NAME) -f build/Dockerfile .

run: build
	docker run -it --name $(CONTAINER_NAME) -e URL=$(URL) -e CONTAINER_NAME=$(CONTAINER_NAME) -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME)

clean:
	docker rm $(CONTAINER_NAME)
	docker rmi $(IMAGE_NAME)
	docker swarm leave -f
