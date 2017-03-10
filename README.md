# echo-xorm

My own toy example with

- HTTP server: [labstack-echo](https://gihtub.com/labstack/echo)
- Database-driver: [go-sqlite3](https://github.com/mattn/go-sqlite3)
- ORM: [xorm](https://github.com/go-xorm/xorm)
- Authorization: [JSON Web Tokens](https://github.com/dgrijalva/jwt-go)


# Installation
## Prerequisites

```bash
go get -u github.com/golang/dep
```

## Application
```bash
go get -u github.com/corvinusz/echo-xorm
dep ensure
```

## Database

Currently using *sqlite3*-database, located at '/tmp/echo-xorm.sqlite.db'

## Vendoring
Used [github.com/golang/dep](https://github.com/golang/dep)

TODO

## Testing
- Test Framework: [Goconvey](https://github.com/smartystreets/goconvey)
- HTTP-Client: [Go-resty](https://github.com/go-resty/resty)

#License

MIT
