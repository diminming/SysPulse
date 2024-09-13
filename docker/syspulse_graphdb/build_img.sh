#!/usr/bin/bash

ENDPOINT="tcp://localhost:8529"
USER="root"
DATABASE="insight"
OUTPUT="/workspace/SysPulse/docker/syspulse_graphdb/insight"

DOCKER_IMAGE_NAME="syspulse_graphdb"

arangodump \
  --server.endpoint ${ENDPOINT} \
  --include-system-collections true \
  --server.username ${USER} \
  --server.database ${DATABASE} \
  --overwrite true \
  --output-directory ${OUTPUT}

docker build -t ${DOCKER_IMAGE_NAME} .