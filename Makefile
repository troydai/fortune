TAG=$(shell git describe --tags)

run:
	@ go run cmd/main.go

image:
	@ docker build . -t echoserver:$(TAG)

docker: image
	@ docker run --rm -p 8080:8080 echoserver

up:
	@ docker-compose up -d

down:
	@ docker-compose down --volumes

# deploy to kubernetes
kup:
	@ kubectl apply -f ./k8s/deployment-dev.yaml

kdown:
	@ kubectl delete -f ./k8s/deployment-dev.yaml

kportal:
	@ kubectl exec -n alice -it portal -- bash
