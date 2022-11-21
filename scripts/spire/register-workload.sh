#!/bin/bash

kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://fortune.troydai.com/ns/spire/sa/spire-agent \
    -selector k8s_sat:cluster:fortune-cluster \
    -selector k8s_sat:agent_ns:spire \
    -selector k8s_sat:agent_sa:spire-agent \
    -node

kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://fortune.troydai.com/ns/fortune/app/front \
    -parentID spiffe://fortune.troydai.com/ns/spire/sa/spire-agent \
    -selector k8s:ns:fortune \
    -selector k8s:pod-label:app.kubernetes.io/name:front 

kubectl exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://fortune.troydai.com/ns/fortune/pod/portal \
    -parentID spiffe://fortune.troydai.com/ns/spire/sa/spire-agent \
    -selector k8s:ns:fortune \
    -selector k8s:pod-name:portal

