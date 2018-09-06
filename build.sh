#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build

docker build -t hicicd .

docker tag hicicd docker.vpclub.cn/openshift/hicicd

docker push docker.vpclub.cn/openshift/hicicd

oc delete imagestream hicicd -n hidevopsio

oc apply -f hicicd.yaml -n  hidevopsio

