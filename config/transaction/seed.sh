#!/bin/bash

kubectl apply -f seed.yaml
jobs=($(kubectl get jobs.batch -o name))
job=${jobs[0]}

echo "waiting for seeding"
kubectl wait --for=condition=complete $job
kubectl delete jobs.batch seed
kubectl delete pvc init-transaction-claim
kubectl delete pvc init-product-claim
kubectl patch pv init-transaction -p '{"spec":{"claimRef": null}}'
kubectl patch pv init-product -p '{"spec":{"claimRef": null}}'
