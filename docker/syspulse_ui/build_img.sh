#!/usr/bin/bash

DOCKER_IMAGE_NAME=syspulse_ui
CURR_DIR=$(cd `dirname $0`;pwd)

cd $CURR_DIR/../../facade/

npm run build:prod

docker build --build-arg http_proxy=http://host.docker.internal:10809 --build-arg https_proxy=http://host.docker.internal:10809 -t ${DOCKER_IMAGE_NAME} .
docker save -o ${CURR_DIR}/${DOCKER_IMAGE_NAME}.tar ${DOCKER_IMAGE_NAME}