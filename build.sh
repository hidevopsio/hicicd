#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build

docker build -t hicicd .

docker tag hicicd docker-registry-default.app.vpclub.io/hidevopsio/hicicd1

docker login -p $(oc whoami -t) -u unused docker-registry-default.app.vpclub.io

docker push docker-registry-default.app.vpclub.io/hidevopsio/hicicd1


