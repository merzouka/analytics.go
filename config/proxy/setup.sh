#!/bin/bash

service="proxy"

kubectl apply -f namespace.yaml
kubectl config set-context --current --namespace $service
