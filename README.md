Login server

Requires Docker version at least 20.10

Backend:    Go

Frontend:   ReactJS

Database:   Postgres


Passwords are stored as Argon2ID hashes. Usernames/Hashes are stored in database, as are sessions. Sessions are stored in browser as cookies and expire after 1 minute.

How to run:

1.  Append 'serverIP mylogin' to etc/hosts
2.  Run docker compose up in the root folder. It will pull postgres and build the images for you
3.  Docker attach to the server container to access the CLI. Use command 'adduser' to add users. CTRL-P, CTRL-Q is the escape sequence to detach from a docker container without killing it.
4.  Fronted webapp is hosted on port 3000. Access with http://mylogin:3000