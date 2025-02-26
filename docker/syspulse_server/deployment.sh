#!/usr/bin/bash

CURR_DIR=$(cd `dirname $0`;pwd)
ENV=$1

source $CURR_DIR/build_img.sh $ENV

if [ $? -eq 0 ]; then
  echo "build docker image success!!!"
else
  echo "build failed..."
fi

scp $CURR_DIR/syspulse_server.${ENV}.tar syspulse_server:/tmp/syspulse_server.tar

ssh syspulse_server 'bash -s' < $CURR_DIR/re-deployment.sh
