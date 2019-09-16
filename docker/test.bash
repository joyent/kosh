#!/usr/bin/env bash

: ${PREFIX:=$USER}
: ${NAME:="kosh"}

: ${BUILDNUMBER:=0}

: ${BUILDER:=${USER}}
BUILDER=$(echo "${BUILDER}" | sed 's/\//_/g' | sed 's/-/_/g')

: ${LABEL:="latest"}
LABEL=$(echo "${LABEL}" | sed 's/\//_/g')

IMAGE_NAME="${PREFIX}/${NAME}:${LABEL}"

docker build \
	-t ${IMAGE_NAME} \
	--file Dockerfile . \
&& \
docker run \
	--rm \
	--name ${BUILDER}_${BUILDNUMBER} \
	${IMAGE_NAME}

