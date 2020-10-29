ACCESS_TOKEN ?= eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI1Mzk1MGEwNS0xOWRiLTExZWItYjg5Ni0wMjQyYWMxMjAwMDIiLCJleHAiOjE2MDM5OTY3OTYsImlhdCI6MTYwMzk5MzE5NiwiaXNzIjoiYXV0aF9zZXJ2aWNlIiwibmJmIjoxNjAzOTkzMTk2fQ.1sxkGtps_4wXN2VO-2tOQKBjphj-x7KX_jn5fkL02gA

.PHONEY: launch_1 launch_10 launch_100 launch_1000

launch_1:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 1 \
		-c 1 \
		-d 10s \
		http://localhost:9999/questionnaires?prefix=m ;\

launch_10:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 10 \
		-c 10 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\

launch_100:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 10 \
		-c 100 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\

launch_1000:
	docker run --rm --net=host \
		williamyeh/wrk  \
		-H "Authorization: ${ACCESS_TOKEN}" \
		-t 1000 \
		-c 1000 \
		-d 1m \
		http://localhost:9999/questionnaires?prefix=m ;\
