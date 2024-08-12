#!/bin/sh

CLUSTER_NAME=analytics

kind create cluster --name $CLUSTER_NAME
kubectl config use-context kind-$CLUSTER_NAME

services=("customer" "product" "transaction")
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
