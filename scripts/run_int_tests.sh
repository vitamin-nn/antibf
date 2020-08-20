#!/bin/bash

trap "docker-compose -f ./deployments/docker-compose.int-tests.yml down" EXIT
docker-compose -f ./deployments/docker-compose.int-tests.yml run integration-tests