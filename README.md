Login server

Requires Docker version at least 20.10

Backend:    Go

Frontend:   ReactJS

Database:   Postgres


Passwords are stored as Argon2ID hashes. User/Hashes are stored in database, as are sessions. Sessions are stored in browser as cookies and expire after 1 minute.

How to run:

1.  Docker pull the official docker postgres image
2.  Docker image build the go-server and webapp images
3.  Append '<server IP> mylogin' to etc/hosts
4.  Run docker compose up to start the server 
5.  To add users, execute '/server adduser' in the go-server container 
6.  Fronted webapp is hosted on port 3000
