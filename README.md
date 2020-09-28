# This is a no factor "Authentication" service aka YOLO Auth

It doesn't provide security, it provides a signed token to services that needs to perform e2e testing

## Contributions

Contributions are very welcome. There are lots that can be improved in this project.

## Setup

It contains a private key in config.go that signs the jwt tokens

And since it's mostly zero config, it starts on <http://0.0.0.0:8089>, all incomming on local server.

## Run

```bash
go build
./no-factor-auth
```

## Docker

```bash
docker run -p 8089:8089 lambdaville/no-factor-auth
```
