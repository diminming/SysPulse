#!/usr/bin/bash

CURR_DIR=$(cd `dirname $0`;pwd)

OUTPUT_DIR=${CURR_DIR}/dist
BIN_NAME="syspulse"

if [ ! -d $OUTPUT_DIR ];then
    mkdir -p $OUTPUT_DIR
fi

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags '-extldflags "-static"' -a -o $OUTPUT_DIR/$BIN_NAME .

cp config_staging.yaml $OUTPUT_DIR/config.yaml
cp bin/startup.sh $OUTPUT_DIR/startup.sh