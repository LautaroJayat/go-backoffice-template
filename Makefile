build-server:
	go build cmd/server/main.go

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

clean:
	rm main