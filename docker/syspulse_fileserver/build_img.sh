#!/usr/bin/bash

# 由于minio基础镜像通过环境变量进行设置即可满足需求，
# 因此未制作镜像，该脚本仅是导出minio的最新镜像

DOCKER_IMAGE_NAME="syspulse_fileserver"

docker save -o ${DOCKER_IMAGE_NAME}.tar minio/minio