#!/usr/bin/bash

CURR_DIR=$(cd `dirname $0`;pwd)

OUTPUT_DIR=${CURR_DIR}/dist
BIN_NAME="syspulse"

if [ ! -d $OUTPUT_DIR ];then
    mkdir -p $OUTPUT_DIR
fi

# go clean -modcache

# export LD_LIBRARY_PATH="-L/workspace/tmp/libpcap-1.10.5"
# export CGO_LDFLAGS="-L/workspace/tmp/libpcap-1.10.5 -extldflags '-static'"
# export CGO_CPPFLAGS="-I/workspace/tmp/libpcap-1.10.5"

# go build -v -o $OUTPUT_DIR/$BIN_NAME
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --ldflags "-extldflags \"-static\"" -a -o $OUTPUT_DIR/$BIN_NAME .

# go build -ldflags='-s -w -L /workspace/tmp/libpcap-1.10.5 -extldflags "-static"'

cp config_staging.yaml $OUTPUT_DIR/config.yaml
