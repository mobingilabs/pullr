#!/bin/bash

docker --version
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
mkdir ${HOME}/.kube
cp ./kubeconf.yaml ${HOME}/.kube/config

# 'master' is our development branch
if [[ "$TRAVIS_BRANCH" == "master" ]]; then
    echo "Setting up development environment"
    kubectl config set clusters.mochi.k8s.local.certificate-authority-data ${KUBE_CLUSTER_CERT}
    kubectl config set clusters.mochi.k8s.local.server ${KUBE_SERVER}
    kubectl config set users.mochi.k8s.local.client-certificate-data ${KUBE_CLIENT_CERT}
    kubectl config set users.mochi.k8s.local.client-key-data ${KUBE_CLIENT_KEYDATA}
fi

if [[ "$TRAVIS_BRANCH" == "production" ]] || [[ "$TRAVIS_BRANCH" == "$TRAVIS_TAG" ]]; then
    echo "Setting up production environment"
    kubectl config set clusters.mochi.k8s.local.certificate-authority-data ${KUBE_PROD_CLUSTER_CERT}
    kubectl config set clusters.mochi.k8s.local.server ${KUBE_PROD_SERVER}
    kubectl config set users.mochi.k8s.local.client-certificate-data ${KUBE_PROD_CLIENT_CERT}
    kubectl config set users.mochi.k8s.local.client-key-data ${KUBE_PROD_CLIENT_KEYDATA}
fi

kubectl version
pip install --user awscli
export PATH=${PATH}:${HOME}/.local/bin
export AWS_ACCESS_KEY_ID=${PULLRCI_ACCESS_KEY_ID};
export AWS_SECRET_ACCESS_KEY=${PULLRCI_SECRET_ACCESS_KEY};
aws --version
