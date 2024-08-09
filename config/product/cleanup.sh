#!/bin/bash

kubectl delete secrets --all

kubectl delete deployments.apps --all
kubectl delete statefulsets.apps --all
kubectl delete pvc --all
kubectl delete pv --all
