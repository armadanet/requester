IMAGE_NAME := requester
CONTAINER_NAME := requester
URL := http://3.23.114.100:9898/newTask
NETWORK := spinner-local-network


build:
	cd ../.. && docker build -t $(IMAGE_NAME) -f requester/build/Dockerfile .

run:
	docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock --net $(NETWORK) --name $(IMAGE_NAME) $(IMAGE_NAME)

clean:
	docker rm $(CONTAINER_NAME)
	docker rmi $(IMAGE_NAME)
	docker swarm leave -f
