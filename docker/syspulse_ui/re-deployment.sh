CONTAINER_ID=$(docker ps -a|awk 'NR>1'|awk '/syspulse_ui/ {print $1}')
IMG_ID=$(docker images|awk 'NR>1'|awk '/syspulse_ui/ {print $3}')


if [ -n "$CONTAINER_ID" ]; then
  docker stop $CONTAINER_ID
  docker rm $CONTAINER_ID
else
  echo "No running syspulse_ui container found."
fi

if [ -n "$IMG_ID" ]; then
  docker image rm $IMG_ID
else
  echo "No syspulse_ui image found."
fi

docker load -i syspulse_ui.tar

docker run -d \
--name syspulse_ui \
-p 8080:80 \
syspulse_ui
