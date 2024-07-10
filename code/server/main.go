package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexedwards/argon2id"
	_ "github.com/lib/pq"
)

func listUsers(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.name, &user.hash)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s %s\n", user.name, user.hash)
	}
}

func listSessions(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM sessions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err := rows.Scan(&session.user, &session.ip, &session.id, &session.expires)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s %s %s %s\n", session.user, session.ip, session.id, session.expires)
	}
}

func userInDB(db *sql.DB, username string) bool {
	username = strings.ToLower(username)
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)

	var user User
	err := row.Scan(&user.name, &user.hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return true
}

func addUser(db *sql.DB, username string, password string) {
	username = strings.ToLower(username)
	if !userInDB(db, username) {
		hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			log.Fatal(err)
		}
		db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hash)
	} else {
		fmt.Printf("%s already in table\n", username) // sql injection vulnerable
	}
}

func dropUser(db *sql.DB, username string) {
	username = strings.ToLower(username)
	if userInDB(db, username) {
		db.Exec("DELETE FROM users WHERE username = $1", username)
		db.Exec("DELETE FROM sessions WHERE username = $1", username)
	} else {
		fmt.Printf("%s does not exist\n", username)
	}
}

func connectToDB() *sql.DB {
	connStr := "postgresql://server:pass@postgres/server?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	args := os.Args
	n := len(args)

	if n > 1 {
		switch {
		case args[1] == "start":
			if n == 2 {
				db := connectToDB()
				Run(db)
			} else {
				fmt.Println("Format: start")
			}
		case args[1] == "users":
			if n == 2 {
				db := connectToDB()
				listUsers(db)
			} else {
				fmt.Println("Format: users")
			}
		case args[1] == "sessions":
			if n == 2 {
				db := connectToDB()
				listSessions(db)
			} else {
				fmt.Println("Format: sessions")
			}
		case args[1] == "adduser":
			if n == 4 {
				db := connectToDB()
				addUser(db, args[2], args[3])
			} else {
				fmt.Println("Format: adduser <username> <password>")
			}
		case args[1] == "dropuser":
			if n == 3 {
				db := connectToDB()
				dropUser(db, args[2])
			} else {
				fmt.Println("Format: dropUser <username>")
			}
		case args[1] == "help":
			if n == 2 {
				fmt.Printf("start - Start the server\nusers - List users in database\nsessions - List sessions in database\nadduser - Add a new user to database\ndropuser - Drop a user from database\n")
			} else {
				fmt.Println("Format: help")
			}
		default:
			fmt.Println("Invalid command")
		}
	} else {
		fmt.Println("Use argument 'start' to start the server or 'help' for a list of commands.")
	}
}
