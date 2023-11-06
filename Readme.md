# Go backoffice template

The present project is a template for a http server that implements a simple crud to handle backoffice database interactions.

As an example, it implements cruds for products and users but it can be extended to include any kind of entity. 

The project uses PostgreSQL for persistence, and Redis as a message broker.

I hope you find this didactic example useful :) 

## Run tests:

For this project we decided to perform tests using an actual Redis and Postgresql connection, more like a Behaviour/acceptance tests, or an integration test.

To try them out you can use our provided docker-compose files and make scripts:

```bash
make run-external

make test

make terminate-external
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
curl --resolve "backoffice.example:80:$(minikube ip)" -i http://backoffice.example/backoffice/products/1 -H "X-Decoded-Perms: 2"
```