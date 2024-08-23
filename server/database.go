package main

import (
	"database/sql"
	"log"
	"strings"

	"github.com/alexedwards/argon2id"
	_ "github.com/lib/pq"
)

func RetrieveUser(db *sql.DB, username string) *sql.Row {
	username = strings.ToLower(username)
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)

	return row
}

func ListUsers(db *sql.DB) []User {
	var users []User
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

		users = append(users, user)
	}

	return users
}

func ListSessions(db *sql.DB) []Session {
	var sessions []Session
	rows, err := db.Query("SELECT * FROM sessions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err := rows.Scan(&session.user, &session.id, &session.expires)
		if err != nil {
			log.Fatal(err)
		}

		sessions = append(sessions, session)
	}

	return sessions
}

func UserInDB(db *sql.DB, username string) bool {
	row := RetrieveUser(db, username)

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

func AddUser(db *sql.DB, username string, password string) bool {
	username = strings.ToLower(username)
	if !UserInDB(db, username) {
		hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			log.Fatal(err)
		}
		db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hash)
		return true
	} else {
		return false
	}
}

func DropUser(db *sql.DB, username string) bool {
	username = strings.ToLower(username)
	if UserInDB(db, username) {
		db.Exec("DELETE FROM users WHERE username = $1", username)
		db.Exec("DELETE FROM sessions WHERE username = $1", username)
		return true
	} else {
		return false
	}
}

func ConnectToDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
