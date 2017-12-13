The `auth.crt` and `auth.key` files were generated using the following command:

```bash
$ openssl req -newkey rsa:4096 -nodes -sha256 -keyout auth.key -x509 -days 365 -out auth.crt
```
