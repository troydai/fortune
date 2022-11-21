TAG=$(shell git describe --tags)
CLUSTER=kind-fortune

.PHONY: build
build:
	@ mkdir -p artifacts
	@ CGO_ENALBE=0 GOOS=linux go build -v -o artifacts/front cmd/front/main.go
	@ CGO_ENALBE=0 GOOS=linux go build -v -o artifacts/datastore cmd/datastore/main.go

run:
	@ go run cmd/main.go

up: build-portal-image build-front-image build-datastore-image
	@ kubectl cluster-info --context kind-$(CLUSTER)
	@ kind load docker-image fortune-datastore:dev fortune-front:dev fortune-portal:dev --name $(CLUSTER)
	@ kubectl apply -f ./k8s/deployment-dev.yaml

down:
	@ kubectl delete -f ./k8s/deployment-dev.yaml

portal:
	@ kubectl exec -n fortune -it portal -- bash

setup-spire:
	@ kind load docker-image \
		gcr.io/spiffe-io/wait-for-it:latest \
		gcr.io/spiffe-io/spire-agent:1.5.0 \
		--name $(CLUSTER)
	@ kubectl apply -f ./k8s/spire/spire-namespace.yaml
	@ kubectl apply \
		-f ./k8s/spire/server-account.yaml \
		-f ./k8s/spire/spire-bundle-configmap.yaml \
		-f ./k8s/spire/server-cluster-role.yaml
	@ kubectl apply \
		-f ./k8s/spire/server-configmap.yaml \
		-f ./k8s/spire/server-statefulset.yaml \
		-f ./k8s/spire/server-service.yaml
	@ kubectl apply \
		-f ./k8s/spire/agent-account.yaml \
		-f ./k8s/spire/agent-cluster-role.yaml
	@ kubectl apply \
		-f ./k8s/spire/agent-configmap.yaml \
		-f ./k8s/spire/agent-daemonset.yaml

# https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/#create-the-cluster
setup-kind:
	@ kind create cluster --name $(CLUSTER) --config=./kind/config.yaml
	# @ cilium install

teardown-kind:
	@ kind delete cluster --name $(CLUSTER)

# https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/#install-the-cilium-cli
cilium-cli:
	@ curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum}
	@ sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
	@ sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
	@ rm cilium-linux-amd64.tar.gz{,.sha256sum}

.PHONY: build-portal-image
build-portal-image:
	@ docker build . --target portal -t fortune-portal:dev

.PHONY: build-front-image
build-front-image:
	@ docker build . --target front -t fortune-front:dev

.PHONY: build-datastore-image
build-database-image:
	@ docker build . --target datastore -t fortune-datastore:dev

