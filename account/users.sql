-- IF user/ role 'golang' not existed yet, you may create it or use your preferred user
-- just make sure to make appropriate change on the script as well ass app config
-- in case you meed to create user 'golang' with password 'golang' you can execute below code 
-- command :
CREATE USER golang with encrypted password 'golang';

-- CREATE DATABASE 'golangtest'
CREATE DATABASE "golangtest" WITH owner golang 
    encoding "UTF8" 
    lc_collate="en_US.UTF-8" 
    lc_ctype="en_US.UTF-8" 
    template 'template0';


-- CREATE SEQUENCE for TABLE users 'id'
CREATE SEQUENCE users_id_seq;

-- Drop table
-- DROP TABLE public.users;

CREATE TABLE public.users (
	id int NOT null DEFAULT nextval('users_id_seq'),
	firstname varchar(30) NOT NULL,
	lastname varchar(30) NULL,
	email varchar(75) NOT NULL,
	passkey varchar(100) NOT NULL,
	CONSTRAINT users_email_un UNIQUE (email),
	CONSTRAINT users_pk PRIMARY KEY (id)
);

-- Permissions

ALTER TABLE public.users OWNER TO golang;
GRANT ALL ON TABLE public.users TO golang;
