# Deployment Analytics Dashboard

[![Build Status](https://travis-ci.org/soprasteria/dad.svg?branch=master)](https://travis-ci.org/soprasteria/dad)
[![Docker Automated buil](https://img.shields.io/docker/automated/soprasteria/dad.svg)](https://hub.docker.com/r/soprasteria/dad/builds/)
[![Go Report Card](https://goreportcard.com/badge/github.com/soprasteria/dad)](https://goreportcard.com/report/github.com/soprasteria/dad)

## Development

Tools and dependencies:
* Golang 1.7
  * [govendor](https://github.com/kardianos/govendor)
* NodeJS 7.2.0
  * npm
  * [gulp](https://github.com/gulpjs/gulp)
* Docker

Get the dependencies:

```sh
npm install
govendor sync
```

Run a MongoDB database:

```sh
docker run --name mongo -p 27017:27017 -v /data/mongo:/data/db -d mongo
```

DAD requires a LDAP configuration. You can write a `~/.dad.toml` file with the following settings:

```toml
env = "prod"

[server]
mongo-addr = "localhost:27017"

[auth]
jwt-secret = "enter a unique pepper here"

[ldap]
address = ""
baseDN = ""
bindDN = ""
bindPassword = ""
searchFilter = ""

[ldap.attr]
username = ""
firstname = ""
lastname = ""
realname = ""
email = ""
```

You can see all the available settings with:

```sh
go run main.go serve --help
```

**Note:** DAD allows three methods for the configuration:

* Use a config file as described above
* Use environment variables (`server.mongo.addr` becomes `DAD_SERVER_MONGO_ADDR`)
* Use CLI parameters (`--server.mongo.addr`)

Run DAD in dev mode, with live reload, with the command:

```sh
gulp
```

You can then browse http://localhost:8080/

## Production

You can generate the binaries with:

```sh
npm run dist
```

The relevant files are `./dad` (the backend) and `./src/dist/` (frontend).

You can generate an archive of these files with:

```sh
npm run archive
```

## License

See the [LICENSE](./LICENSE) file.
