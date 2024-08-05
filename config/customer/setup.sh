#!/bin/bash

kind create cluster --name analytics
kubectl apply -f namespace.yaml
kubectl config set-context --current --namespace customer

db_service="db-service"
db_password=$(echo -n $(cat .db))
db_url="postgres://docker:$db_password/$db_service.customer.svc.cluster.local:5432/customerdb"
kubectl create secret generic db-secret --from-literal=db-password=$db_password --from-literal=db-url=$db_url
