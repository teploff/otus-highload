.PHONY: lint build swagger_spec

build:
	go build cmd/messenger/main.go

lint:
	golangci-lint run --enable-all ;\

swagger_spec:
	 $(GOPATH)/bin/swag init --generalInfo="./cmd/messenger/main.go" --dir="./" --output="./api/swagger" ;\
