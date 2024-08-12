#!/bin/bash

service="transaction"

kubectl apply -f namespace.yaml
kubectl config set-context --current --namespace $service

get_service_url () {
    echo "$1.$service.svc.cluster.local"
}

db_service="database"
db_name="$service""db"
db_password=$(echo -n $(cat .db))
kubectl create secret generic db-secret --from-literal=db-password=$db_password

kubectl create secret generic service-secrets \
    --from-literal=db-url="postgresql://docker:$db_password@$(get_service_url "$db_service"):5432/$db_name" \
    --from-literal=cache-url="$(get_service_url "cache"):6379" \
    --from-literal=cache-password=$(echo -n $(cat .cache))

kubectl apply -f init.yaml
