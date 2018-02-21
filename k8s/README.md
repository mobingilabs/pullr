# Running on kubernetes

Create these secrets and configmaps before creating any k8s resources. Feel free
to change the secret values as you wish.

```
kubectl create secret generic pullr \
    --from-file=pullr.crt=certs/pullr.crt \
    --from-file=pullr.key=certs/pullr.key \
    --from-literal=mongorootuser=root \
    --from-literal=mongorootpass=rootpass \
    --from-literal=mongopass=pullrpass \
    --from-literal='regpass=admin' \
    --from-literal='githubid=<put your github oauth client id here>' \
    --from-literal='githubsecret=<put your github oauth secret here>'
    
    
kubectl create configmap pullr \
    --from-file=conf \
    --from-literal=mongouser=pullr \
    --from-literal=reguser=admin
```
