#!/bin/bash

echo "[E2E PreTest] Prepare for the environment"

#record current dir and root dir for the convenience
pretestBaseDir=$(pwd)
cd ../..
workspaceBaseDir=$(pwd)
cd ${pretestBaseDir}

# prepare necessary environment for e2e test
# 0 remove all old logs
cd "${workspaceBaseDir}/test/log"
rm -f *.log

#exit if any command encountered error
set -e

echo "[E2E PreTest] Check existence of minikube"
# 1.check whether minikube exists
hasMinikube=$(command -v minikube | wc -l)
if [ $hasMinikube == 0 ]
then
    echo "[E2E PreTest] minikube not found, installing it"
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    sudo install minikube-linux-amd64 /usr/local/bin/minikube
else
    echo "[E2E PreTest] minikube found at $(command -v minikube)"
fi

# 2.totally fresh minikube start
echo "[E2E PreTest] start minikube environment"
minikube delete
#minikube start
minikube start --image-mirror-country='cn'








