#!/bin/bash

# login to ecr
`aws ecr get-login --region ap-northeast-1`

IMAGE_TAG="latest"

if [[ "$TRAVIS_BRANCH" == "$TRAVIS_TAG" ]]; then
    IMAGE_TAG=${TRAVIS_TAG}
fi

IMAGE="${NAME}:${IMAGE_TAG}"
docker tag ${IMAGE} ${ECR_REPO_URI}/${IMAGE}
docker images
docker push ${ECR_REPO_URI}/${IMAGE}
