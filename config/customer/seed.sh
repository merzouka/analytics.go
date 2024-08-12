#!/bin/bash

kubectl apply -f seed.yaml
jobs=($(kubectl get jobs.batch -o name))
job=${jobs[0]}

echo "waiting for seeding"
kubectl wait --for=condition=complete $job
kubectl delete jobs.batch seed
kubectl delete pvc init-$service-claim
kubectl patch pv init-$service -p '{"spec":{"claimRef": null}}'
