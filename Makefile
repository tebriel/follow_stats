.PHONY: build run

build:
	docker build -t cmoultrie/follow_stats .

run:
	docker run --rm -it \
		-e CONSUMER_KEY \
		-e CONSUMER_SECRET \
		-e ACCESS_TOKEN \
		-e ACCESS_TOKEN_SECRET \
		-e ES_URL=http://frodux.in:9200 \
		-l SERVICE_NAME=follow_stats \
		-l SERVICE_TAGS=http \
		--dns 10.77.2.12 \
		-p 9090:8080 \
		cmoultrie/follow_stats \
		app -v=2 -logtostderr=true
