# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=isaac
VERSION=`cat VERSION`
DOCKER_IMAGE_NAME=iconloop/$(BINARY_NAME)
GIT_REVISION=`git rev-parse --short HEAD`
GREP_TEMPORARY_FILE=`grep -E  '(log$ |data$ |/.db$)' `

# Build
#===================
all: build test
build:doc
	# Build frontend
	$(MAKE) -C frontend
  # Build backend
	$(GOGET) -d -v
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	./clean.sh
	rm -f $(BINARY_NAME)
	rm -rf frontend/build
	rm -rf frontend/node_modules
	rm -f  data/*.*
	find . -name "*.log" -exec rm -f {} \;

install:
	$(GOBUILD) install
doc:
	rm -rf docs/
	swag init
test:
	# Frontend test
	$(MAKE) -C frontend test
	# Backend test
	go test -count 1 motherbear/backend/...
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

# Cross build
#===================
build-linux:
	# Build frontend
	cd frontend && $(MAKE)
    # Build backend  for linux
	$(GOGET) -d -v
	GOOS=linux CC=gcc GOARCH=amd64 $(GOBUILD)  -a -installsuffix cgo -o $(BINARY_NAME) -v .

# Docker operation
#===================
docker-build:build
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
docker-build-dev:build
	docker build -t $(DOCKER_IMAGE_NAME):$(GIT_REVISION) .
	docker tag $(DOCKER_IMAGE_NAME):$(GIT_REVISION) $(DOCKER_IMAGE_NAME):dev
docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest
docker-push-dev:
	docker push $(DOCKER_IMAGE_NAME):$(GIT_REVISION)
	docker push $(DOCKER_IMAGE_NAME):dev
