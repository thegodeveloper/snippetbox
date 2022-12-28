# Snippetbox Web Application

Web application with Go and PostgreSQL database. This repository is my study notes of the *[Alex Edwards book - Let's Go](https://lets-go.alexedwards.net/)*.

<img src="img/snippetbox.png" alt="Snippetbox" style="float: left; margin-right: 10px;" />

## PostgreSQL Database

### Start PostgreSQL Service

```
sudo service postgresql start
```

### Create PostgreSQL Database

Option 1:

```
sudo -u postgres -i

psql -U postgres
psql (13.3)
Type "help" for help.

postgres=# create database snippetbox;
CREATE DATABASE
postgres=# create user hachiko with encrypted password 'nirvana';
CREATE ROLE
postgres=# grant all privileges on database snippetbox to hachiko;
GRANT
postgres=#
```

Option 2:

```
sudo -u postgres psql
psql (13.9 (Debian 13.9-0+deb11u1))
Type "help" for help.

postgres=#
```

### Connect to snippetbox database

```
\connect snippetbox
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
