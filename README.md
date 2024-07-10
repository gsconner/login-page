Login server

Requires Docker at least 20.10

Backend:    Go

Frontend:   ReactJS

Database:   Postgres


Passwords are stored as Argon2ID hashes. User/Hashes are stored in database, as are sessions. Sessions are stored in browser as cookies and expire after 1 minute.


How to run:

1.  Docker pull the official docker postgres image
2.  
3.  Docker image build the go-server and webapp images
4.  
5.  Append '<server IP> mylogin' to etc/hosts
6.  
7.  Run docker compose up to start the server
8.  
9.  To add users, execute '/server adduser' in the go-server container
10.  
11.  Fronted webapp is hosted on port 3000
