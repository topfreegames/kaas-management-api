.PHONY: all dep build

all: dep build

dep:
	@echo "  >  Making sure go.mod matches the source code"
	go mod tidy -v

build:
	@echo "  >  build"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o build/manager

initial-setup: init-kind init-cluster-api

deploy: init-kind init-tilt

init-kind:
	bash ./scripts/kind.sh kaas-cluster

init-cluster-api:
	clusterctl init --infrastructure docker

wait-cluster-api-resources:
	bash ./scripts/wait-cluster-api.sh

create-cluster:
	clusterctl generate cluster capi-quickstart --flavor development \
   		--kubernetes-version v1.22.0 \
    	--control-plane-machine-count=3 \
    	--worker-machine-count=3 | kubectl apply -f -

init-tilt:
	kind export kubeconfig --name kaas-cluster
	tilt up
