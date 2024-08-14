#!/bin/bash

namespaces=("customer" "transaction" "product" "proxy")

for ns in ${namespaces[@]}; do
    printf "\e[1m                               $ns\e[m\n"
    kubectl config set-context --current --namespace $ns
    kubectl delete secrets --all
    kubectl delete jobs.batch --all
    kubectl delete statefulsets.apps --all
    kubectl delete deployments.app --all
    kubectl delete svc --all
    kubectl delete pvc --all
done

printf "\e[1m                               general\e[m\n"
kubectl delete pv --all
