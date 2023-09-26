all: test build push

test:
	go test ./... -v -cover
build: test
	docker build . -t registry.digitalocean.com/ryansheppard/weather
push: build
	docker push registry.digitalocean.com/ryansheppard/weather
