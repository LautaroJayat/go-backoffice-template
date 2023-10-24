build-server:
	go build cmd/server/main.go

run-db:
	sudo docker compose -f DockerCompose.test.yaml up --build

terminate-db:
	sudo docker compose -f DockerCompose.test.yaml down

run-redis:
	sudo docker compose -f DockerCompose.propagation.test.yaml up --build

terminate-redis:
	sudo docker compose -f DockerCompose.propagation.test.yaml down

test-api:
# first run-db
	go test ./api/...

test-propagation:
# first run-redis
	go test ./propagation/...

clean:
	rm main
