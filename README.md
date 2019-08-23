# Go Cache Store

A cache store server with REST interface.

## Start

```sh
cat <<EOT >> .env
PORT=3001
SENTRY_DSN=https://yoursentryendpoint.io
EOT

docker build -t cache-server .
```

```sh
docker swarm init
docker stack deploy -c docker-compose.yml <stack-name>
```

## Usage

+ *GET* `/cache/:key` -> get the registry by the key
+ *POST* `/cache/:key` -> create a new key/value registry
+ *DELETE* `/cache/:key` -> delete the registry stored by the key

+ *GET* `/health` -> health check the server
+ *GET* `/metrics` -> prometheus application metrics
