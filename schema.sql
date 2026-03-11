CREATE TABLE users(
    id serial PRIMARY KEY,
    username varchar(100) UNIQUE NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password varchar(255) NOT NULL,
    created_at timestamp DEFAULT now()
);