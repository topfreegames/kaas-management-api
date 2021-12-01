.PHONY: all dep build test lint fix

all: dep build

dep:
	@echo "  >  Making sure go.mod matches the source code"
	go mod tidy -v

build:
	@echo "  >  build"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o build/manager

test:
	@echo "> Running tests"
	go test ./...

lint:
	@echo " > Running golangci-lint"
	golangci-lint run

fix:
	go fmt ./...

# Development environment targets
setup-dev-env: create-kind init-tilt

create-kind:
	bash ./scripts/kind.sh kaas-cluster

init-tilt:
	kind export kubeconfig --name kaas-cluster
	tilt up

tilt-ci:
	kind export kubeconfig --name kaas-cluster
	tilt ci

# Tilt targets
apply-capi-dependencies:
	kubectl apply -f ./scripts/assets/crds/dependencies

wait-capi-dependencies-resources:
	bash ./scripts/wait-controllers.sh dependencies

apply-capi:
	kubectl apply -f ./scripts/assets/crds/cluster-api

wait-capi-resources:
	bash ./scripts/wait-controllers.sh cluster-api

apply-test-clusters:
	kubectl apply -f ./scripts/assets/crs