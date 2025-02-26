CONTAINER_ID=$(docker ps -a|awk 'NR>1'|awk '/syspulse_server/ {print $1}')
IMG_ID=$(docker images|awk 'NR>1'|awk '/syspulse_server/ {print $3}')

docker stop $CONTAINER_ID
docker rm $CONTAINER_ID
docker image rm $IMG_ID
docker load -i /tmp/syspulse_server.tar

docker run -d \
--name syspulse_server \
-p 24160:24160 \
-p 24162:24162 \
-e "FS_ACCESS_KEY=g0ipu3vpDkZhTOVmIBkd" \
-e "FS_SECRET_KEY=gN6o0V5sO2fzZQ1lUe2gcddM6G6sViDBUF6o9NGb" \
syspulse_server
