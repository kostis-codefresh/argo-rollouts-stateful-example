# Using Argo Rollouts with stateful services

Example with Argo Rollouts with two services that use a queue (RabbitMQ in our case).
We want to do a blue/green deployment without affecting the production queue.

## Local development

This is the workflow when you need to change the source code.

### Run RabbitMQ locally

Install Docker 

```
docker run --name some-rabbit -p 5672:5672 -p 5673:5673 -p 15672:15672 rabbitmq:3-management
```
You can visit the dashboard at `http://localhost:15672`.
Username and password are `guest:guest`

### Run the worker manually

Install [GoLang](https://go.dev/) locally

```
cd source/worker
APP_VERSION=1.2  go run .
```
You can now access the worker at `http://localhost:8080`

Click the button to test RabbitMQ connection - the worker sends
messages to itself

### Run the tester/producer locally

```
cd source/tester
go run .
```
You can now access the producer at `http://localhost:7000`

Click the *Production* button to send messages to the worker


### Run all services at the same time

Install [Docker compose](https://docs.docker.com/compose/) (no need for local GoLang installation)

```
cd src
docker compose up
```

And now you can use 

* `http://localhost:15672` for RabbitMQ
* `http://localhost:8000` for Stable worker
* `http://localhost:9000` for Preview Worker


### Deploy on Kubernetes

```
cd manifests/stateful-rollout
kubectl apply -f . 
kubectl port-forward svc/rabbitmq 15672:15672
kubectl port-forward svc/my-plain-frontend-service 9000:8080 -n plain
```

And now you can use 

* `http://localhost:15672` for RabbitMQ
* `http://localhost:8000` for Stable worker

or launch the file `example.html` in your browser found at the root 
of the repository.

## Start a rollout

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


