: ${IMAGE_NAME:=asssaf/crickithat:latest}
docker run --rm -it --privileged --device /dev/i2c-1 "$IMAGE_NAME" $*
