#!/bin/bash

kind create cluster --name analytics
kubectl apply -f namespace.yaml
kubectl config set-context --current --namespace customer

db_service="db-service"
db_password=$(echo -n $(cat .db))
db_url="postgres://docker:$db_password/$db_service.customer.svc.cluster.local:5432/customerdb"
kubectl create secret generic db-secret --from-literal=db-password=$db_password --from-literal=db-url=$db_url

kubectl apply -f init.yaml

pods=($(kubectl get pods -o name))
pod=${pods[0]}
kubectl wait --for=condition=Ready $pod

init_folder="/docker-entrypoint-initdb.d"
kubectl exec $pod -- psql -U docker -d customerdb -h localhost -a -f "$init_folder/01-create-tables.sql"
kubectl exec $pod -- psql -U docker -d customerdb -h localhost -a -f "$init_folder/02-populate-tables.sql"

kubectl delete jobs.batch seed
kubectl delete pvc init-claim
kubectl delete pv init 
kubectl delete pvc data-claim
