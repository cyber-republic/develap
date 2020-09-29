#!/bin/bash

while getopts ":p:l:d:n:h" opt; do
  case $opt in
    h)
      echo "Usage:"
      echo "    tools/build_dockerimages.sh -h                  Display this help message."
      echo "    tools/build_dockerimages.sh -p [yes|no]         Push built images to docker registry. You will need access to push to https://hub.docker.com/u/cyberrepublic to use this option."
      echo "    tools/build_dockerimages.sh -l [yes|no]         Whether to tag the docker image as 'latest'"
      echo "    tools/build_dockerimages.sh -d dockerhub_namespace         Pass the dockerhub namespace to use to build images for"
      echo "    tools/build_dockerimages.sh -n [all|mainchain|did-sidechain|eth-sidechain]         Build a specific node"
      exit 0
    ;;
    p) DOCKER_PUSH="$OPTARG"
    ;;
    l) DOCKER_PUSH_LATEST="$OPTARG"
    ;;
    d) DOCKER_NAMESPACE="$OPTARG"
    ;;
    n) DOCKER_IMAGE_TO_BUILD="$OPTARG"
    ;;
    \? )
      echo "Invalid Option: -$OPTARG" 1>&2
      exit 1
    ;;
  esac
done

if [ -z "${DOCKER_NAMESPACE}" ]
then
  echo "You must pass in a docker namespace with '-d' flag"
  exit 1
fi


function build_docker_image {
    REPO_URL="${1}"
    REPO_BRANCH="${2}"
    DOCKER_NAME="${3}"
    DOCKER_TAG="${4}"

    CURRENT_DIR=$(pwd)
    DOCKER_IMAGE="${DOCKER_NAMESPACE}/ela-${DOCKER_NAME}"

    docker build -t "${DOCKER_IMAGE}:latest" -f "${CURRENT_DIR}/docker/${DOCKER_NAME}.Dockerfile" .
    if [ "${DOCKER_PUSH}" == "yes" ]
    then
      if [ "${DOCKER_PUSH_LATEST}" == "yes" ]
      then
        docker push "${DOCKER_IMAGE}:latest"
      fi
      docker tag "${DOCKER_IMAGE}:latest" "${DOCKER_IMAGE}:${DOCKER_TAG}"
      docker push "${DOCKER_IMAGE}:${DOCKER_TAG}"
    fi
}

if [ "${DOCKER_IMAGE_TO_BUILD}" == "mainchain" ] || [ "${DOCKER_IMAGE_TO_BUILD}" == "all" ]
then 
    build_docker_image "github.com/elastos/Elastos.ELA" "v0.6.0" "mainchain" "v0.6.0"
fi
if [ "${DOCKER_IMAGE_TO_BUILD}" == "did-sidechain" ] || [ "${DOCKER_IMAGE_TO_BUILD}" == "all" ]
then 
    build_docker_image "github.com/elastos/Elastos.ELA.SideChain.ID" "v0.2.0" "did-sidechain" "v0.2.0"
fi
if [ "${DOCKER_IMAGE_TO_BUILD}" == "eth-sidechain" ] || [ "${DOCKER_IMAGE_TO_BUILD}" == "all" ]
then 
    build_docker_image "github.com/elastos/Elastos.ELA.SideChain.ETH" "v0.1.0" "eth-sidechain" "v0.1.0"
fi