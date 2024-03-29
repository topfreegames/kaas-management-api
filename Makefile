.PHONY: all dep build test lint fix

all: fix lint test dep build

dep:
	@echo "  >  Making sure go.mod matches the source code"
	go mod tidy -v
	go install github.com/swaggo/swag/cmd/swag@latest

build: dep build-docs
	@echo "  >  build"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o build/manager

build-docs:
	@echo " > Running swaggo"
	 swag init -g internal/server/server.go

test:
	@echo "> Running tests"
	go test ./...

lint:
	@echo " > Running golangci-lint"
	golangci-lint run

fix:
	@echo " > Running go fmt"
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