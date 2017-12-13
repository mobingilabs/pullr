#!/bin/bash

DEPLOY=0

if [[ "$TRAVIS_BRANCH" == "master" ]] || [[ "$TRAVIS_BRANCH" == "production" ]]; then
    DEPLOY=1
fi

IMAGE_TAG=`echo ${TRAVIS_COMMIT}`

if [[ "$TRAVIS_BRANCH" == "$TRAVIS_TAG" ]]; then
    IMAGE_TAG=${TRAVIS_TAG}
    DEPLOY=1
fi

if [[ $DEPLOY -eq 0 ]]; then
    echo "Deployment not supported"
    exit 0
fi

# login to ecr
export AWS_ACCESS_KEY_ID=`echo ${PULLRCI_ACCESS_KEY_ID}`
export AWS_SECRET_ACCESS_KEY=`echo ${PULLRCI_SECRET_ACCESS_KEY}`
`aws ecr get-login --no-include-email --region ap-northeast-1`

if [ $? -ne 0 ]; then
    exit 1
fi

# push image to ecr
IMAGE="${NAME}:${IMAGE_TAG}"
echo "image = ${IMAGE}"
docker tag ${IMAGE} ${ECR_REPO_URI}/${IMAGE}
docker images
docker push ${ECR_REPO_URI}/${IMAGE}

if [ $? -ne 0 ]; then
    exit 1
fi

# update kubernetes deployment with the new image
kubectl set image deployment ${NAME} ${NAME}=${ECR_REPO_URI}/${IMAGE}
