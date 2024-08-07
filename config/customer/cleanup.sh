#!/bin/bash

kubectl delete secrets --all

kubectl delete pvc --all
kubectl delete -f init.yaml
kubectl delete -f cache.yaml
kubectl delete -f db.yaml
kubectl delete -f customer.yaml
