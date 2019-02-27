# Go Cache Store

A cache store server with REST interface.

## Start

```bash
$ docker build -t <image-name> .
```

```bash
$ docker run -p 3001:3001 <image-name>
```

## Usage

+ *GET* `/cache/:key` -> get the registry by the key
+ *POST* `/cache/:key` -> create a new key/value registry
+ *DELETE* `/cache/:key` -> delete the registry stored by the key

+ *GET* `/check` -> health check the server
