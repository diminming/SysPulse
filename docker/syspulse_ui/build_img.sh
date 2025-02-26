#!/usr/bin/bash

DOCKER_IMAGE_NAME=syspulse_ui
CURR_DIR=$(cd `dirname $0`;pwd)

# 按照目标环境构建镜像，备选项有：
# staging： 验证环境
# cmcc-xj： 新疆移动生产
# cic： 中华保险

# default is staging
ENV=$1

cd $CURR_DIR/../../facade/

npm run build:${ENV}

docker build --no-cache --build-arg ENVIRONMENT=$ENV -t ${DOCKER_IMAGE_NAME} .

docker save -o ${CURR_DIR}/${DOCKER_IMAGE_NAME}.${ENV}.tar ${DOCKER_IMAGE_NAME}
