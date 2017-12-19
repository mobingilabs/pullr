## Overview

tbd

## Getting started

Install the following prerequisites first:

- [docker](https://docs.docker.com/engine/installation/)
- [docker-compose](https://docs.docker.com/compose/install/)

Set the following environment variables as well:

- `APISERVER_ACCESS_KEY_ID`
- `APISERVER_SECRET_ACCESS_KEY`

> Tested only on Linux (Ubuntu)

```bash
# go to admin to add superuser
$ docker exec -it mongodb mongo admin

> use admin;
> db.createUser({user: "root", pwd: "rootpass", roles: ["root"]});

# connect using superuser
$ docker exec -it mongodb mongo -u root -p rootpass --authenticationDatabase admin

# create pullr user for tokenserver
> db.createUser({user: "pullr", pwd: "pullrpass", roles: ["readWrite"]});

> use pullr;
> db.createCollection("users");

# insert admin/admin user to collection
# the bcrypt hash for the password admin was generated using `htpasswd -nB admin`
> db.users.insert({"username": "admin", "password": "$2y$05$oBNfJkZ4rMd6PjrRHq3FdeZXezfBzWqWsZuJ7v0ePpdUFCVNaOv52"});
> db.users.find({});

# then run locally
$ make up

# test docker login
# valid user/password combinations
#   admin / admin
$ docker login localhost:5000

# when done, run
$ make down
```
