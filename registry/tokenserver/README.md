The file `auth_config.yml` is copied from [this file](https://github.com/cesanta/docker_auth/blob/master/examples/simple.yml).

## Mongodb notes

```bash
# go to admin to add superuser
$ docker exec -it mongodb mongo admin

> use admin;
> db.createUser({ user: "root", pwd: "rootpass", roles: ["root"] });

# connect using superuser
$ docker exec -it mongodb mongo -u root -p rootpass --authenticationDatabase admin

> use pullr;
> db.createCollection("users");

# insert admin/admin user to collection
> db.users.insert({"username": "admin", "password": "$2y$05$oBNfJkZ4rMd6PjrRHq3FdeZXezfBzWqWsZuJ7v0ePpdUFCVNaOv52"});
> db.users.find({});
```
