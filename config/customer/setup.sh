#!/bin/bash

kind create cluster --name analytics
kubectl apply -f namespace.yaml
kubectl config set-context --current --namespace customer

get_service_url () {
    echo "$1.customer.svc.cluster.local"
}

db_service="db-service"
db_password=$(echo -n $(cat .db))
db_url="postgres://docker:$db_password/$db_service.customer.svc.cluster.local:5432/customerdb"
kubectl create secret generic db-secret --from-literal=db-password=$db_password --from-literal=db-url=$db_url

kubectl create secret generic customer-secrets \
    --from-literal=db-url="postgresql://docker:$db_password@$(get_service_url "database"):5432" \
    --from-literal=cache-url="$(get_service_url "cache"):6379" \
    --from-literal=cache-password=$(echo -n $(cat .cache))

kubectl apply -f init.yaml

kubectl delete jobs.batch seed
kubectl delete pvc init-claim
kubectl patch pv init -p '{"spec":{"claimRef": null}}'

# database set up
kubectl apply -f db.yaml

# cache set up
kubectl apply -f cache.yaml

# customer set up
kubectl apply -f customer.yaml
