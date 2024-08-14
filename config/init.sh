#!/bin/sh

CLUSTER_NAME=analytics

kind create cluster --name $CLUSTER_NAME --config cluster.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml 
kubectl apply -f ingress.yaml
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
kubectl config use-context kind-$CLUSTER_NAME

services=("customer" "product" "transaction" "proxy")
for svc in ${services[@]}; do
    printf "\e[1m                               $svc\e[m\n"
    cd ./$svc/
    chmod +x *.sh
    bash setup.sh
    if [ "$1" == "--seed" ]; then
        bash seed.sh
    fi
    bash deploy.sh
    cd ..
done
