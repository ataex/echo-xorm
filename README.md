# echo-xorm

My own toy example with

- HTTP server: [labstack-echo](https://gihtub.com/labstack/echo)
- Database-driver: [go-sqlite3](https://github.com/mattn/go-sqlite3)
- ORM: [xorm](https://github.com/go-xorm/xorm)
- Authorization: [JSON Web Tokens](https://github.com/dgrijalva/jwt-go)


# Installation
## Prerequisites
Currently using *dep* as a dependencies manager.

```bash
go get -u github.com/golang/dep/cmd/dep
```

## Installation process
```bash
go get -d github.com/corvinusz/echo-xorm
cd $GOPATH/src/github.com/corvinusz/echo-xorm/
dep ensure
go install
```

## Vendoring
Implemented via [github.com/golang/dep](https://github.com/golang/dep)
To update vendor dependencies use
```bash
cd $GOPATH/src/github.com/corvinusz/echo-xorm/
dep ensure -update
```

## Configuration
Config file located as

`$GOPATH/src/github.com/corvinusz/echo-xorm/deploy/echo-xorm-config.toml`
Feel free to change.

By default application expect to find the configuration file as:

`/usr/local/etc/echo-xorm-config.toml`

You can point path to config file it with '-config' flag

## Database
Currently using:
- *sqlite3*-database, located as '/tmp/echo-xorm.sqlite.db' (change it in config)
- ORM [xorm](https://github.com/go-xorm/xorm)

# Application Run
```bash
echo-xorm -h # shows application flags
echo-xorm -config=$GOPATH/src/github.com/corvinusz/echo-xorm/deploy/echo-xorm-config.toml # runs app with default cfg
```

To health check
```bash
curl http://localhost:11111/version
```

It should return something like

`{"result":"OK","version":"0.0.1","server_time":1501286982}`

Of course you feel free to use any of applicable applications such as:
- [insomnia](https://insomnia.rest/)
- [postman](https://www.getpostman.com/)
- etc..

or you browser plugin like:
- [chromerestclient](https://advancedrestclient.com/)
- [restclient](https://addons.mozilla.org/ru/firefox/addon/restclient/)
- etc...

# Testing
Implemented in BDD-style with:
- Test Framework: [gomega](https://github.com/onsi/gomega)
- HTTP-Client: [Go-resty](https://github.com/go-resty/resty)

```bash
cd $GOPATH/src/github.com/corvinusz/echo-xorm/bddtests
go test
```

Test parameters defined in file:

`$GOPATH/src/github.com/corvinusz/echo-xorm/bddtests/test-config/echo-xorm-test-config.toml`

#License

MIT
