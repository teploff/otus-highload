.PHONY: init migrate launch_backend_wmn down_backend_wmn launch_backend_wsn down_backend_wsn wrk_1 wrk_2 down

init:
	docker-compose up -d storage_master storage_slave_1 storage_slave_2;\

migrate:
	cd ../server/tools/migrator && \
		go build main.go && \
		mv main ../../../cluster && \
		cd - ;\
	./main --dir="../server/migrations" --dsn="user:password@tcp(localhost:3306)/social-network" "up" ;\
	rm main ;\

launch_backend_wmn:
	docker-compose up --build -d backend_with_master_node ;\

down_backend_wmn:
	docker rm -f backend_with_master_node ;\

launch_backend_wsn:
	docker-compose up --build -d backend_with_slave_node ;\

down_backend_wsn:
	docker rm -f backend_with_slave_node ;\

wrk_1:
	echo ${ACCESS_TOKEN} ;\
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 10 \
		-c 10 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\

wrk_2:
	echo ${ACCESS_TOKEN} ;\
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 10 \
		-c 10 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=j ;\

down:
	docker-compose down ;\
