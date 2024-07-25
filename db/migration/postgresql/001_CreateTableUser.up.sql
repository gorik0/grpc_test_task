CREATE TABLE users
(
    id SERIAL,
	email VARCHAR(150) NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_email_key UNIQUE (email)
);