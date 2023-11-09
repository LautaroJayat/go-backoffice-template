build-server:
	go build cmd/server/main.go

build-auth-proxy:
	go build cmd/auth_proxy/main.go

build-docker-backoffice:
	docker build -t backoffice -f Dockerfile.app .

build-docker-auth-proxy:
	docker build -t auth_proxy -f Dockerfile.proxy .

build-token-generator:
	go build -o ./genToken proxy/tokenGenerator/main.go

run-external:
	sudo docker compose -f DockerCompose.test.yaml up --build

terminate-external:
	sudo docker compose -f DockerCompose.test.yaml down

run-redis:
	sudo docker compose -f DockerCompose.propagation.test.yaml up --build

terminate-redis:
	sudo docker compose -f DockerCompose.propagation.test.yaml down

test-api:
# first run-external
	go test ./api/...

test-propagation:
# first run-redis or run-external
	go test ./propagation/...

test:
# first run-external
	go test ./...

test-proxy:
	echo "not implemented yet"
	echo "not implemented yet"
	echo "not implemented yet"


gen-test-keys:
	rm -drf proxy/.tmp
	mkdir proxy/.tmp
	openssl genrsa -out proxy/.tmp/private.pem 2048
	openssl rsa -in proxy/.tmp/private.pem -pubout > proxy/.tmp/public.pem

clean:
	rm -drf proxy/.tmp
	rm -f main

