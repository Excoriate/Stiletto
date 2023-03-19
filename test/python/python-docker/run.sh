#!/usr/bin/env bash

set -e
set -o pipefail

function docker_build() {
    docker build -t ${CONTAINER_NAME} .
}

function docker_run_detached() {
    docker run -d -p 8080:8080 --name ${CONTAINER_NAME} ${CONTAINER_NAME}
}

function run_stiletto_ci(){
  echo "Running stiletto ci"
}

declare CONTAINER_NAME="python-fastapi-docker-stiletto"

while [[ $# -gt 0 ]]
do
  key="$1"
  case $key in
      -build)
      docker_build
      shift
      ;;
      -run)
      docker_run_detached
      shift
      ;;
      -ci)
      run_stiletto_ci
      shift
      ;;
      *)
      echo "Invalid option: $key"
      exit 1
      ;;
  esac
done
