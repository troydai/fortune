#!/bin/bash

kubectl delete pod -n fortune --force portal
kubectl apply -f ./k8s/portal.yaml