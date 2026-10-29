# Chirpy
A place where you can chirp to the Internet.

(~~an alternative to X.~~)

## Motivation
This project is the product of a guided [Boot.dev](https://www.boot.dev/courses/learn-http-servers-golang) project to understand and exercise building an HTTP server in Go.

### Goals
- Build an HTTP server from (almost) scratch
- Communicate with client via a RESTful API using Go's standard [http package](https://pkg.go.dev/net/http)
- Use Postgres as a database
- Implement a authentication and authorization system using [jwt package](https://pkg.go.dev/github.com/golang-jwt/jwt)
- Build a webhook and provide an API key to authenticate with
- This documentation you're reading

## Prerequisites
- Go 1.22+
- PostgresSQL
- [Goose](http://github.com/pressly/goose) for database migrations

## Getting Started
### 1. Clone the repo
```bash
git clone https://codeberg.org/kedarui/chirpy.git
```

### 2. Set up PostgresSQL

#### Installation
**macOS**
```bash
brew install postgres@14
```

**Linux**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Ensure the installation worked:
```bash
psql --version
# ex: psql (PostgreSQL) 14.24 (Ubuntu 14.24-0ubuntu0.22.04.1)
```

**For Linux** (optional), you can set up a password for your postgres account:
```bash
sudo passwd postgres
```

#### Start up PostgresSQL service
**macOS**
```bash
brew service start postgresql@14
```

**Linux**
```bash
sudo service postgresql start
```

#### Connect to server
**macOS**
```bash
psql postgres
```

**Linux**
```bash
sudo -u postgres psql 
# as admin (assuming you are...), you can connect as the postgres account
```

You will then see this:
```
postgres=#
```

#### Create a new database
```sql
CREATE DATABASE chirpy;
```

#### Connect to database
```
\c chirpy
```

Prompt should change to:
```
chirpy=#
```

**For Linux**, set the user password:
```sql
ALTER USER postgres WITH PASSWORD 'somethingyouwillremember'
# this is used for your db connection string
```

To exit the `psql` shell, type `exit` or `\q`.

#### Connection String
Speaking of which, you should know what a connection string looks like. They all follow this format:
```bash
protocol://username:password@host:port/database?parameters
```

So for our Postgres use-case:
```bash
postgres://postgres:somethingyouwillremember@localhost:5432/chirpy?sslmode=disable

```

- **protocol** -> `postgres`
- **username** -> `postgres`, remember from earlier?
- **password** -> `somethingyouwillremember`, Only for Linux. I guess because we Linux users are secure like that.
- **host** -> `localhost`, assuming you're using your own machine
- **port** -> `5432`, PostgresSQL default. It's theirs. Made official by the [IANA](https://www.iana.org/assignments/service-names-port-numbers?search=&page=93).
- **database** -> `chirpy` (this project)
- **parameters** (optional) -> `sslmode=disable`, if you're self-hosting and playing around within your local machine, then you can disable this. If you plan on self-hosting and making this public, then you should research SSL/TLS. These are protocols for protecting the data in your database.

### 3. Migrate your goose
If you don't already have [Goose](http://github.com/pressly/goose):
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

The data migrations come with the repo, so within `sql/schema`, run:
```bash
goose <connection_string> up
```

This will get your Postgres database up-to-date.

Recommend checking out Goose's documentation to make migrations easier. Can you run them without typing your connection string every time?

### 4. Environment
Create a `.env` at the root of the project and add the following:
```env
DB_URL="connection_string"
PLATFORM="dev" # for testing
TOKEN_SECRET="yoursecretforauthentication"
```

### 5. Go run
```bash
go run .
```

Now Chirpy *should* be running.
