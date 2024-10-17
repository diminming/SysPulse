#!/usr/bin/bash

DOCKER_IMAGE_NAME="syspulse_graphdb"

ENDPOINT="tcp://localhost:8529"
USER="root"
DATABASE="insight"
OUTPUT="/workspace/SysPulse/docker/syspulse_graphdb/syspulse"

if [ ! -d "${OUTPUT}" ]; then
  mkdir -p ${OUTPUT}
fi

arangodump \
  --server.endpoint ${ENDPOINT} \
  --include-system-collections true \
  --server.username ${USER} \
  --server.database ${DATABASE} \
  --overwrite true \
  --output-directory ${OUTPUT}

docker build -t ${DOCKER_IMAGE_NAME} .
docker save -o ${DOCKER_IMAGE_NAME}.tar ${DOCKER_IMAGE_NAME}