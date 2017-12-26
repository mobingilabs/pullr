## Overview

tbd

## Getting started

Install the following prerequisites first:

- [docker](https://docs.docker.com/engine/installation/)
- [docker-compose](https://docs.docker.com/compose/install/)

Set the following environment variables as well:

- `APISERVER_ACCESS_KEY_ID`
- `APISERVER_SECRET_ACCESS_KEY`

> Tested only on Linux (Ubuntu 16.04)

```bash
# run the stack first so we can connect to db and create users
$ make up

# connect to db as admin and add root superuser
$ docker exec -it mongodb mongo admin

> use admin;
> db.createUser({user: "root", pwd: "rootpass", roles: ["root"]});
> exit

# connect to db using root
$ docker exec -it mongodb mongo -u root -p rootpass --authenticationDatabase admin

# create pullr user for our tokenserver
> use pullr;
> db.createUser({user: "pullr", pwd: "pullrpass", roles: ["readWrite"]});
> db.createCollection("users");

# insert admin/admin user to collection, used for docker login
# the bcrypt hash for the password admin was generated using `htpasswd -nB admin`
> db.users.insert({"username": "admin", "password": "$2y$05$oBNfJkZ4rMd6PjrRHq3FdeZXezfBzWqWsZuJ7v0ePpdUFCVNaOv52"});
> db.users.find({});
> exit

# then rerun the whole stack
$ make down && make up

# test `docker login` using "admin/admin" as credentials
$ docker login localhost:5000

# when done, run
$ make down
```
