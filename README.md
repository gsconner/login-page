# Login Page

* Backend: Go
* Frontend: ReactJS
* Database: Postgres

Requires Docker version at least 20.10

## Description

A login page supported by a Go backend and a webapp created with ReactJS, with a Postgres database. User data is stored in the database in the form of username/hash pairs. The passwords are hashed and salted via Argon2id. When login credentials are entered on the webapp, they are sent to the Go server via an HTTP POST request, which then compares the user's password to the corresponding entry in the database. If the log-in is successful, the server creates a session and sends back a cookie to the webapp, which is used to verify the ongoing session for future requests. The sessions are also stored in the database and expire after 1 minute.

## How to install

To run the server, use docker compose up in the main folder. It will build docker images for each component, and pull a postgres docker image of the correct version, then start the corresponding containers. 
Once the containers are running, use docker attach on the server container to access the CLI for the server. From there you can run commands and add users so you can log in. The website is hosted on http://localhost:3000.