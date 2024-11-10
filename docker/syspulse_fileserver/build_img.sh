#!/usr/bin/bash

DOCKER_IMAGE_NAME="syspulse_fileserver"

docker build -t ${DOCKER_IMAGE_NAME} .
docker save -o ${DOCKER_IMAGE_NAME}.tar ${DOCKER_IMAGE_NAME}