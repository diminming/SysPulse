#!/usr/bin/bash

CURR_DIR=$(cd `dirname $0`;pwd)
nohup ./inisght conf config.yaml > startup.log 2>&1 &
