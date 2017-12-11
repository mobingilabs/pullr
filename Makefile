# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: pre release

release: pre custom 

pre:
	@docker --version; \
	curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl; \
	chmod +x ./kubectl; \
	sudo mv ./kubectl /usr/local/bin/kubectl; \
	mkdir ${HOME}/.kube; \
	cp ./kubeconf.yaml ${HOME}/.kube/config; \
    kubectl config set clusters.mochi.k8s.local.certificate-authority-data ${KUBE_CLUSTER_CERT}; \
    kubectl config set clusters.mochi.k8s.local.server ${KUBE_SERVER}; \
    kubectl config set users.mochi.k8s.local.client-certificate-data ${KUBE_CLIENT_CERT}; \
    kubectl config set users.mochi.k8s.local.client-key-data ${KUBE_CLIENT_KEYDATA}; \
	kubectl version; \
	pip install --user awscli; \
    export PATH=${PATH}:${HOME}/.local/bin; \
    aws --version;
