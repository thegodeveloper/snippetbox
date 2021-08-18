CREATE TABLE public.snippets (
	id serial NOT NULL,
	title varchar NOT NULL,
	"content" text NOT NULL,
	created timestamp NOT NULL,
	expires timestamp NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE public.users (
	id serial NOT NULL,
	name varchar NOT NULL,
	email varchar NOT NULL,
	hashed_password char(60) NOT NULL,
	created timestamp NOT NULL,
	active boolean NOT NULL DEFAULT TRUE
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2018-12-23 17:25:22'
);
