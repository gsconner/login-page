package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

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
				users := ListUsers(db)
				for _, user := range users {
					fmt.Printf("%s %s\n", user.name, user.hash)
				}
			} else {
				fmt.Println("Format: users")
			}
		case args[0] == "sessions":
			if n == 1 {
				sessions := ListSessions(db)
				for _, session := range sessions {
					fmt.Printf("%s %s %s\n", session.user, session.id, session.expires)
				}
			} else {
				fmt.Println("Format: sessions")
			}
		case args[0] == "adduser":
			if n == 3 {
				if !AddUser(db, args[1], args[2]) {
					fmt.Printf("%s already in table\n", args[1])
				}
			} else {
				fmt.Println("Format: adduser <username> <password>")
			}
		case args[0] == "dropuser":
			if n == 2 {
				if !DropUser(db, args[1]) {
					fmt.Printf("%s does not exist\n", args[1])
				}
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
		db := ConnectToDB(args[1])
		go Run(db)
		console(db)
	} else {
		fmt.Printf("Please include a link to the database.")
	}
}
