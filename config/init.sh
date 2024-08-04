#!/bin/sh

CLUSTER_NAME=analytics

kind create cluster --name $CLUSTER_NAME
kubectl config use-context kind-$CLUSTER_NAME
