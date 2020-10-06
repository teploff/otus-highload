.PHONY: lint up down

lint:
	golangci-lint run --enable-all ;\

up:
	docker-compose -f deployments/docker-compose.yml up -d db ;\
	sleep 2 ;\
	docker-compose -f deployments/docker-compose.yml up -d adminer ;\
#	go build -o goose tools/migrator/main.go ;\
#	./goose \
#		--dir="./migrations" \
#		--dsn="user:password@tcp(localhost:3306)/social-network" \
#		up ;\
#	rm goose ;\
#        docker-compose -f deployments/stage/docker-compose.yml up -d --build auth ;\
#        docker image prune -f ;\

down:
	docker-compose -f deployments/docker-compose.yml down ;\