# Chirpy
---
A place where you can chirp to the Internet.

(~~an alternative to X.~~)

## Motivation
---
This project is the product of a guided [Boot.dev](https://www.boot.dev/courses/learn-http-servers-golang) project to understand and exercise building an HTTP server in Go.

### Goals
- Build an HTTP server from (almost) scratch
- Communicate with client via a RESTful API using Go's standard [http package](https://pkg.go.dev/net/http)
- Use Postgres as a database
- Implement a authentication and authorization system using [jwt package](https://pkg.go.dev/github.com/golang-jwt/jwt)
- Build a webhook and provide an API key to authenticate with
- This documentation you're reading

## Prerequisites
---
- Go 1.22+
- PostgresSQL
- [Goose](http://github.com/pressly/goose) for database migrations
