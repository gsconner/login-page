Login server

Requires Docker version at least 20.10

Backend:    Go

Frontend:   ReactJS

Database:   Postgres


Passwords are stored as Argon2ID hashes. Usernames/Hashes are stored in database, as are sessions. Sessions are stored in browser as cookies and expire after 1 minute.

How to run:

1.  Docker pull the official docker postgres image
2.  Docker image build the server and webapp images from the dockerfiles. Make sure you name them 'go-server' and 'webapp' respectively
3.  Append 'serverIP mylogin' to etc/hosts
4.  Run docker compose up (in the docker folder so it uses compose.yaml) to start the server 
5.  Docker attach to the server container to access the CLI. Use command 'adduser' to add users. CTRL-P, CTRL-Q is the escape sequence to detach from a docker container without killing it.
6.  Fronted webapp is hosted on port 3000. Access with http://mylogin:3000