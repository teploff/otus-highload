.PHONY: init down

init:
	docker-compose up -d ;\

migrate:
	cd ../server/tools/migrator && \
		go build main.go && \
		mv main ../../../cluster && \
		cd - ;\
	./main --dir="../server/migrations" --dsn="user:password@tcp(localhost:3306)/social-network" "up" ;\
	rm main ;\

down:
	docker-compose down ;\
