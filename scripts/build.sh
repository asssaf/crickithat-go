: ${IMAGE_NAME:=asssaf/crickithat:latest}
BASE="$(dirname $0)/.."
docker build -t $IMAGE_NAME -f $BASE/docker/Dockerfile $BASE
