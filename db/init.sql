CREATE TABLE users (
    username varchar,
    password varchar
);

CREATE TABLE sessions (
    username varchar,
    sessid varchar,
    expires timestamp
);