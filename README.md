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

### ZooKeeper

| param | type   | default           | description                  | example           |
| :---- | :----- | :---------------- | :--------------------------- | :---------------- |
| ZK_SERVICES    | string | 127.0.0.1:2182    | zookeeper connection address | 127.0.0.1:2182    |
| ZK_USERNAME    | string |                   | zookeeper username           | root              |
| ZK_PASSWD    | string |                   | zookeeper passwd             | 123456            |
| ZK_CONFIG_PATH    | string | /visualization/config| project config file zNode    | /visualization/config |
