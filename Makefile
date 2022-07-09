TAG=$(shell git describe --tags)
CLUSTER=kind-fortune

run:
	@ go run cmd/main.go

up:
	@ kubectl cluster-info --context kind-$(CLUSTER)
	@ docker build . -t fortune-datastore:dev -f docker/datastore/Dockerfile
	@ docker build . -t fortune-front:dev -f docker/front/Dockerfile
	@ docker build . -t fortune-portal:dev -f docker/portal/Dockerfile
	@ kind load docker-image fortune-datastore:dev fortune-front:dev fortune-portal:dev --name $(CLUSTER)
	@ kubectl apply -f ./k8s/deployment-dev.yaml

down:
	@ kubectl delete -f ./k8s/deployment-dev.yaml

portal:
	@ kubectl exec -n fortune -it portal -- bash

setup-kind:
	@ kind create cluster --name $(CLUSTER) --wait 10m

teardown-kind:
	@ kind delete cluster --name $(CLUSTER)
