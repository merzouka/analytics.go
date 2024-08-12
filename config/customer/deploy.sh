#!/bin/bash

# database set up
kubectl apply -f db.yaml

# cache set up
kubectl apply -f cache.yaml

# service set up
kubectl apply -f service.yaml
