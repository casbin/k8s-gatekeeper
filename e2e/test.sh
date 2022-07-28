#!/bin/bash

pretestBaseDir=$(pwd)
cd ..
workspaceBaseDir=$(pwd)

# 3.build webhook as external service
echo "[E2E PreTest] build admission webhook"
cd $workspaceBaseDir
pwd
go mod tidy
go mod vendor

echo "[E2E PreTest] load Model and Policy CRD to k8s"
cd "${workspaceBaseDir}"
kubectl apply -f config/auth.casbin.org_casbinmodels.yaml
kubectl apply -f config/auth.casbin.org_casbinpolicies.yaml

echo "[E2E Test] test Start"
cd $workspaceBaseDir
go test -v --count=1 -run ^TestE2E$ github.com/casbin/k8s-gatekeeper/e2e

