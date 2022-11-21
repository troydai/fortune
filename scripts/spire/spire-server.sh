#!/bin/bash

kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server $@
