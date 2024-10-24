#!/usr/bin/bash

CURR_DIR=$(cd `dirname $0`;pwd)

source $CURR_DIR/build_img.sh

if [ $? -eq 0 ]; then
  echo "build docker image success!!!"
else
  echo "build failed..."
fi

scp $CURR_DIR/syspulse_ui.tar syspulse_server:/home/admin

ssh syspulse_server 'bash -s' < $CURR_DIR/re-deployment.sh
