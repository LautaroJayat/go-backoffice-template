# Go backoffice template

The present project is a template for a http server that implements a simple crud to handle backoffice database interactions, a reverse proxy that handles JWT signature validation, and a kubernetes deployment. 

As an example, it implements cruds for products and users but it can be extended to include any kind of entity. 

The project uses PostgreSQL for persistence, and Redis as a message broker.

I hope you find this didactic example useful :) 

## Architecture

![Static Diagram](static_diagram.png)

As we can see in the diagram, the project has been conceived to be deployed in a kubernetes cluster.

The flow is the following:
1. An ingress gets requests from the internet.
2. It forwards the requests to the auth proxy deployment.
3. The proxy decodes and validates the JWT present in the Auth header.
4. If the JWT is not valid or is not present, it delivers an early response.
5. If everything is OK, it changes the request adding some headers.
6. Finally it forwards the request to the main application.
7. The application is responsible for checking permissions and interacting with the DB (Postgres) and the message broker (Redis).

## Run tests

For this project we decided to perform tests using an actual Redis and Postgresql connection, more like a Behavior/acceptance tests, or an integration test.

To try them out you can use our provided docker-compose files and make scripts:

```bash
# this will create RSA keys for the auth proxy
make gen-test-keys

# this will run redis and postgres
make run-external

# this will run all tests
make test

# this will terminate all containers we created
make terminate-external

# and this will clean the generated RSA keys
make clean
```

## Minikube deployment

If you want to try this out in minikube you can follow these steps.

### 1. Enable addons

```bash
minikube addons enable helm-tiller
minikube addons enable ingress 
```

### 2. Build and push docker image into minikube's docker registry

```bash
eval $(minikube docker-env)
docker build -t backoffice
docker push backoffice
```

### 3. Deploy Redis and Postgres

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install -f _k8s/helm_redis/values.yaml redis bitnami/redis
helm install -f _k8s/helm_postgres/values.yaml postgres bitnami/postgresql
```

### 4. Deploy the application
```bash
minikube kubectl -- create -f ./_k8s
```

### 5. Test endpoints

After this, you will be able to reach our application through the ingress.

The simplest way to interact is by using `curl`.

Remember that you will need to provide value for `X-Decoded-Role` header so the auth middleware allows you to reach the actual endpoint.

See [role package](./roles/roles.go) and [middleware package](./api/http/middleware/middleware.go) 

```bash
curl \
--resolve "backoffice.example:80:$(minikube ip)" \
-i http://backoffice.example/backoffice/products/ \
-H "X-Decoded-Perms: 1"
```

## The Monkey
![The Monkey](the_monke.jpg)