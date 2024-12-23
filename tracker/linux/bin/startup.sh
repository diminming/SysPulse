#!/usr/bin/bash

if [ "$(id -u)" -ne 0 ]; then
    echo "The script must be run as root. Please use sudo or run as the root user"
    exit 1
fi

CURR_DIR=$(cd `dirname $0`;pwd)

SYSPULSE=$CURR_DIR/syspulse
CONF_FILE=$CURR_DIR/config.yaml
LOG_PATH=/tmp/syspulse.log

nohup $SYSPULSE --conf $CONF_FILE > $LOG_PATH 2>&1 &