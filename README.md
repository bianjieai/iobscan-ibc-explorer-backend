# iobscan-ibc-explorer-backend


## Run

First make

```bash
make build
```

then start with

```bash
./iobscan-ibc-explorer-backend start
```

or

```bash
./iobscan-ibc-explorer-backend start test -c configFilePath
```

## Run with docker

You can run application with docker.

```bash
docker build -t iobscan-ibc-explorer-backend .
```

then

```bash
docker run --name iobscan-ibc-explorer-backend -p 8080:8080 iobscan-ibc-explorer-backend
```

## env params
- CONFIG_FILE_PATH: `option` `string` config file path
