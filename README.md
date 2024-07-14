Login server

Requires Docker version at least 20.10

Backend:    Go

Frontend:   ReactJS

Database:   Postgres


Passwords are stored as Argon2ID hashes. Usernames/Hashes are stored in database, as are sessions. Sessions are stored in browser as cookies and expire after 1 minute.

To run the server, use docker compose up in the main folder. It will build docker images for each component, and pull a postgres docker image of the correct version, then start the corresponding containers. 
Once the containers are running, use docker attach on the server container to access the CLI for the server. From there you can run commands and add users so you can log in. The website is hosted on http://localhost:3000.