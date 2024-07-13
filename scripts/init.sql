CREATE TABLE users (
    username varchar,
    password varchar
);

CREATE TABLE sessions (
    username varchar,
    sessid varchar,
    ip varchar,
    expires timestamp
);