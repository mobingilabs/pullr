## Overview

tbd

## Getting started

```bash
# install the following prerequisites first:
#   docker
#   docker-compose

# you also need you have these environment variables (todo: make optional)
#   APISERVER_ACCESS_KEY_ID
#   APISERVER_SECRET_ACCESS_KEY

# then run locally
$ make up

# test docker login
$ docker login localhost:5000
Username: admin
Password: badmin
Login Succeeded

# when done, run
$ make down
```
