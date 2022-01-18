-- users table
CREATE TABLE public.users (
	id serial NOT NULL,
	firstname varchar(30) NOT NULL,
	lastname varchar(30) NULL,
	email varchar(75) NOT NULL,
	passkey varchar(100) NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_un UNIQUE (email)
);

