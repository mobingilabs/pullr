#!/bin/bash

# login to ecr
export AWS_ACCESS_KEY_ID=`echo ${PULLRCI_ACCESS_KEY_ID}`;
export AWS_SECRET_ACCESS_KEY=`echo ${PULLRCI_SECRET_ACCESS_KEY}`;
`aws ecr get-login --region ap-northeast-1`

IMAGE_TAG=`echo ${TRAVIS_COMMIT}`

if [[ "$TRAVIS_BRANCH" == "$TRAVIS_TAG" ]]; then
    IMAGE_TAG=${TRAVIS_TAG}
fi

IMAGE="${NAME}:${IMAGE_TAG}"
echo "image = ${IMAGE}"
docker tag ${IMAGE} ${ECR_REPO_URI}/${IMAGE}
docker images
docker push ${ECR_REPO_URI}/${IMAGE}
