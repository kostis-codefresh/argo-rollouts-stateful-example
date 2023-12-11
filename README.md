# Using Argo Rollouts with stateful services

Example with Argo Rollouts with two services that use a queue (RabbitMQ in our case).
We want to do a blue/green deployment without affecting the production queue.


## Run RabbitMQ locally

Install Docker 

```
docker run --name some-rabbit -p 5672:5672 -p 5673:5673 -p 15672:15672 rabbitmq:3-management
```
You can visit the dashboard at `http://localhost:15672`.
Username and password are `guest:guest`

## Run only the backend manually

Install [GoLang](https://go.dev/) locally

```
cd source/interest
APP_VERSION=1.2  go run .
```
You can now access the backend at `http://localhost:8080`



## Run all services at the same time

Install [Docker compose](https://docs.docker.com/compose/) (no need for local GoLang installation)

```
cd src
docker compose up
```

And now you can use the same URLs as above to access the services.

## Run on Kubernetes as deployments

```
cd manifests/plain
kubectl create ns plain
kubectl apply -f . -n plain
kubectl port-forward svc/my-plain-backend-service 8000:8080 -n plain
kubectl port-forward svc/my-plain-frontend-service 9000:8080 -n plain
```

You can now access the backend at `http://localhost:8000` and the backend at `http://localhost:9000`

## Run on Kubernetes as Rollouts (modern app)

```
cd manifests/modern
kubectl create ns modern
kubectl apply -f . -n modern
kubectl port-forward svc/backend-active 8000:8080 -n modern
kubectl port-forward svc/backend-preview 8050:8080 -n modern

kubectl port-forward svc/frontend-active 9000:8080 -n modern
kubectl port-forward svc/frontend-preview 9050:8080 -n modern
```

You can now access the backend at `http://localhost:8000` (old) and `http://localhost:8050` (new)
and the backend at `http://localhost:9000` (old) and `http://localhost:9050` (new)

To see what the rollouts are doing

```
kubectl-argo-rollouts get rollout my-frontend -n modern
kubectl-argo-rollouts get rollout my-backend -n modern
```


