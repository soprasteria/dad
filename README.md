# Deployment Analytics Dashboard

[![Build Status](https://travis-ci.org/soprasteria/dad.svg?branch=master)](https://travis-ci.org/soprasteria/dad)
[![Docker Automated buil](https://img.shields.io/docker/automated/soprasteria/dad.svg)](https://hub.docker.com/r/soprasteria/dad/builds/)
[![Go Report Card](https://goreportcard.com/badge/github.com/soprasteria/dad)](https://goreportcard.com/report/github.com/soprasteria/dad)
[![Code Coverage](https://codecov.io/gh/soprasteria/dad/branch/master/graph/badge.svg)](https://codecov.io/gh/soprasteria/dad)
[![Dependencies Status](https://david-dm.org/soprasteria/dad/status.png)](https://david-dm.org/soprasteria/dad)
[![Dev Dependencies Status](https://david-dm.org/soprasteria/dad/dev-status.png)](https://david-dm.org/soprasteria/dad?type=dev)

## Development

Tools and dependencies:
* [Golang 1.7](https://golang.org/)
  * [govendor](https://github.com/kardianos/govendor)
* [NodeJS 8](https://nodejs.org/en/)
* [Docker](https://www.docker.com/)

**Don't forget to add $GOPATH/bin to your $PATH**

## Get the dependencies

```sh
npm install
govendor sync
```

## Run a MongoDB database

```sh
docker run --name mongo -p 27017:27017 -v /data/mongo:/data/db -d mongo
```

## Specify the server configuration

DAD requires a LDAP configuration. You can write a `~/.dad.toml` file with the following settings:

```toml
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

[docktor]
addr = "http://<DocktorUrl>/#!/"
user = "<DocktorUsername>"
password = "<DocktorPassword>"

[tasks]
recurrence = "@every 20m"

```

You can see all the available settings with:

```sh
go run main.go serve --help
```

**Note:** DAD allows three methods for the configuration:

* Use a config file as described above
* Use environment variables (`server.mongo.addr` becomes `DAD_SERVER_MONGO_ADDR`)
* Use CLI parameters (`--server.mongo.addr`)

## Run the project

Run DAD in dev mode, with live reload, with the command:

```sh
npm start
```

You can then browse [http://localhost:8080/](http://localhost:8080/)

## Populate the database with data

You can use the shell scripts in [bin/](./bin). The scripts use CSV files. The format of these files is described in the shell scripts.

## Production

You can generate the binaries with:

```sh
npm run dist
```

The relevant files are in the `dist` folder.

## Run deployment analytics job

In order to run the routine to analyses which functional services of projects are deployed or not, you can POST a request to the endpoint API `/api/admin/jobs/deployment-indicators` with an admin account.

Every project with Docktor URL will be updated with data from Docktor.

This kind of routine is also executed at regular time (default to 23:00 everyday, can be overridden with `--tasks-recurrence` option)

## License

See the [LICENSE](./LICENSE) file.
