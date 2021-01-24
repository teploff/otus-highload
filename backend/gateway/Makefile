.PHONY: lint build swag_mix

build:
	go build cmd/gateway/main.go

lint:
	golangci-lint run --enable-all ;\

swag_mix:
	docker pull quay.io/goswagger/swagger ;\
	alias swagger="docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger" ;\
	swagger version ;\
	swagger mixin -o api/swagger/swagger.json \
		api/swagger/auth.swagger.json \
		api/swagger/social.swagger.json ;\
