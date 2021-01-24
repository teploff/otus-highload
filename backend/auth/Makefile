.PHONY: lint build swagger_spec

build:
	go build cmd/auth/main.go

lint:
	golangci-lint run --enable-all ;\

swagger_spec:
	 $(GOPATH)/bin/swag init --generalInfo="./cmd/auth/main.go" --dir="./" --output="./api/swagger" ;\
