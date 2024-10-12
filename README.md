# Snippetbox Web Application

Web application with Go and PostgreSQL database. This repository is my study notes of the *[Alex Edwards book - Let's Go](https://lets-go.alexedwards.net/)*.

<img src="img/snippetbox.png" alt="Snippetbox" style="float: left; margin-right: 10px;" />

## PostgreSQL Database

### Start PostgreSQL Database

- Docker image called `cool_black`
- Start the container

### Connect to PostgreSQL Database

```shell
psql postgres -h localhost -U postgres                                                                                                                                       ─╯
Password for user postgres: password
psql (17.0, server 15.3 (Debian 15.3-1.pgdg110+1))
Type "help" for help.

postgres=#
```

### Create PostgreSQL Database

Option 1:

```
postgres=# create database snippetbox;
CREATE DATABASE
postgres=# create user hachiko with encrypted password 'nirvana';
CREATE ROLE
postgres=# grant all privileges on database snippetbox to hachiko;
GRANT
postgres=#
```

### Connect to snippetbox database

```
psql --host=localhost --dbname=greenlight --username=hachiko
Password for user hachiko: nirvana
psql (17.0, server 15.3 (Debian 15.3-1.pgdg110+1))
Type "help" for help.

greenlight=>
```

### Connect using an environment variable

```
psql $GREENLIGHT_DB_DSN
psql (13.9 (Debian 13.9-0+deb11u1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.

greenlight=>
```

### Create snippets table

```
CREATE TABLE public.snippets (
	id serial NOT NULL,
	title varchar NOT NULL,
	"content" text NOT NULL,
	created timestamp NOT NULL,
	expires timestamp NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);
```

### Create users table

```
CREATE TABLE public.users (
	id serial NOT NULL,
	name varchar NOT NULL,
	email varchar NOT NULL,
	hashed_password char(60) NOT NULL,
	created timestamp NOT NULL,
	active boolean NOT NULL DEFAULT TRUE
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

### How to insert test information

```
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + interval '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + interval '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP + interval '365 days'
);
```

## Generating a self-signed TLS

```
mkdir tls && cd tls

go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Run the application

```
go run snippetbox.godeveloper.net/cmd/web
```

## Build

```
go build -o /tmp/snippetbox snippetbox.godeveloper.net/cmd/web

/tmp/snippetbox
```
