#!/bin/bash

service="customer"

kubectl apply -f namespace.yaml
kubectl config set-context --current --namespace $service

get_service_url () {
    echo "$1.$service.svc.cluster.local"
}

db_service="database"
db_name="$servicedb"
db_password=$(echo -n $(cat .db))
kubectl create secret generic db-secret --from-literal=db-password=$db_password

kubectl create secret generic service-secrets \
    --from-literal=db-url="postgresql://docker:$db_password@$(get_service_url "$db_service"):5432/$db_name" \
    --from-literal=cache-url="$(get_service_url "cache"):6379" \
    --from-literal=cache-password=$(echo -n $(cat .cache))

kubectl apply -f init.yaml
pods=($(kubectl get pods -o name))
pod=${pods[0]}

echo "waiting for seeding"
kubectl wait --for=condition=Ready $pod
kubectl delete jobs.batch seed
kubectl delete pvc init-claim
kubectl patch pv init -p '{"spec":{"claimRef": null}}'

# database set up
kubectl apply -f db.yaml

# cache set up
kubectl apply -f cache.yaml

# service set up
kubectl apply -f service.yaml
