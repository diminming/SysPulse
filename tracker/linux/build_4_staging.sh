#!/usr/bin/bash

CURR_DIR=$(cd `dirname $0`;pwd)

OUTPUT_DIR=${CURR_DIR}/dist
BIN_NAME="syspulse"

if [ ! -d $OUTPUT_DIR ];then
    mkdir -p $OUTPUT_DIR
fi

CC=musl-gcc go build -ldflags '-extldflags "-static -lpcap -ldbus-1 -lsystemd"' -o $OUTPUT_DIR/$BIN_NAME

cp config_staging.yaml $OUTPUT_DIR/config.yaml


