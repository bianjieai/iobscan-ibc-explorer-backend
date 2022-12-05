# iobscan-ibc-openapi


## Run

First make

```bash
make build
```

then start with

```bash
./iobscan-ibc-openapi start
```

or

```bash
./iobscan-ibc-openapi start test -c configFilePath
```

## Run with docker

You can run application with docker.

```bash
docker build -t iobscan-ibc-openapi .
```

then

```bash
docker run --name iobscan-ibc-openapi -p 8080:8080 iobscan-ibc-openapi
```

## env params
- CONFIG_FILE_PATH: `option` `string` config file path
