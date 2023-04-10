#!/usr/bin/env bash

set -e -o pipefail

declare CONTAINER_NAME="example-web-container"
declare IMAGE_NAME="example-web"

function removeContainerIfExist(){
  if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    docker rm -f $CONTAINER_NAME
  fi
}

function build(){
  docker build -t $IMAGE_NAME . && \
  docker run -d --name $CONTAINER_NAME -p 8080:8080 $IMAGE_NAME
}

function testContainer(){
  curl "http://localhost:8080/" -H "accept: application/json" | jq
}

function main(){
  removeContainerIfExist
  build
  sleep 3 # to ensure that the container started
  testContainer
}

main
