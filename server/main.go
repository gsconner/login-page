package main

import (
	"bufio"
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

func connectToDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func console(db *sql.DB) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		line := scanner.Text()

		args := strings.Split(line, " ")
		n := len(args)

		switch {
		case args[0] == "users":
			if n == 1 {
				listUsers(db)
			} else {
				fmt.Println("Format: users")
			}
		case args[0] == "sessions":
			if n == 1 {
				listSessions(db)
			} else {
				fmt.Println("Format: sessions")
			}
		case args[0] == "adduser":
			if n == 3 {
				addUser(db, args[1], args[2])
			} else {
				fmt.Println("Format: adduser <username> <password>")
			}
		case args[0] == "dropuser":
			if n == 2 {
				dropUser(db, args[1])
			} else {
				fmt.Println("Format: dropUser <username>")
			}
		case args[0] == "help":
			if n == 1 {
				fmt.Printf("users - List users in database\nsessions - List sessions in database\nadduser <username> <password> - Add a new user to database\ndropuser <username> - Drop a user from database\n")
			} else {
				fmt.Println("Format: help")
			}
		default:
			if args[0] != "" {
				fmt.Println("Invalid command. Type 'help' for a list of commands.")
			}
		}
	}
}

func main() {
	args := os.Args
	if len(args) >= 2 {
		db := connectToDB(args[1])
		go Run(db)
		console(db)
	} else {
		fmt.Printf("Please include a link to the database.")
	}

}
