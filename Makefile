GOSRC = $(GOPATH)/src
DOCKER_EXEC = docker run -v $(PWD):/go/src/github.com/mikemackintosh/wonka -it wonka

docker:
	docker build -t wonka .

build:
	$(DOCKER_EXEC) go run main.go

test:
	echo $(GOSRC)
	$(DOCKER_EXEC) go test ./...

useradd:
	$(DOCKER_EXEC) go build -o /go/src/github.com/mikemackintosh/wonka/bin/useradd cmd/useradd/useradd.go
	$(DOCKER_EXEC) /go/src/github.com/mikemackintosh/wonka/bin/useradd -l

spec:
	$(DOCKER_EXEC) sh -c "cd testing && bundle exec rake spec"
