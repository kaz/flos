# FLOS: Log Observing System

FLOS is server observation and administration suites for swift incident responses.
FLOS can introduce into Linux/Darwin/Windows instances.

## Features

- Automatic network forming (instance detection)
- Run command on multiple instances
- Alive monitoring for any services
- File access logging (with fanotify, only for linux)
- Log aggregation and flexible quering by SQL
- Automatic file backup and incremental snapshots
- Safe from hijacking (strong encrypted communication)

## Build

See `Makefile`.

## Run

FLOS can run with single binary.
Specify `-d` to run in background.

```sh
$ flos -d
```

## Management

FLOS uses `10239/tcp` and `10239/udp` to communicate with manage server.

Management serverr is served on other repository: [flos-hortus](https://github.com/kaz/flos-hortus)
