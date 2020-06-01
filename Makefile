GOSRC = $(GOPATH)/src
PKG_PATH = $(shell git rev-parse --show-toplevel | sed -e "s|^\($(GOPATH)\)/src/||")
DOCKER_EXEC = docker run -v $(PWD):/go/src/$(PKG_PATH) -it wonka

env:
	@echo "PKG_PATH=$(PKG_PATH)"

docker:
	docker build -t wonka .

build:
	$(DOCKER_EXEC) go run main.go

test:
	echo $(GOSRC)
	$(DOCKER_EXEC) go test ./...

useradd:
	$(DOCKER_EXEC) go build -o /go/src/$(PKG_PATH)/bin/useradd cmd/useradd/useradd.go
	$(DOCKER_EXEC) /go/src/$(PKG_PATH)/bin/useradd -l

spec:
	$(DOCKER_EXEC) sh -c "cd testing && bundle exec rake spec"
