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

init-dependencies:
	kubectl apply -f ./scripts/assets/crds/dependencies

wait-dependencies-resources:
	bash ./scripts/wait-controllers.sh dependencies

init-cluster-api:
	kubectl apply -f ./scripts/assets/crds/cluster-api

wait-cluster-api-resources:
	bash ./scripts/wait-controllers.sh cluster-api

create-clusters:
	kubectl apply -f ./scripts/assets/crs

init-tilt:
	kind export kubeconfig --name kaas-cluster
	tilt up

tilt-ci:
	kind export kubeconfig --name kaas-cluster
	tilt ci