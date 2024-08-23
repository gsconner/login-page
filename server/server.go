package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	_ "github.com/lib/pq"
)

type User struct {
	name string
	hash string
}

type Session struct {
	user    string
	id      string
	expires time.Time
}

type Arg struct {
	name  string
	value string
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func startSession(db *sql.DB, username string, sessID string) {
	_, err := db.Exec("INSERT INTO sessions (username, sessid, expires) VALUES ($1, $2, $3)", username, sessID, time.Now().UTC().Add(time.Minute).Format(time.DateTime))
	if err != nil {
		log.Fatal(err)
	}
}

func genID() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

func createCookie(username string, sessID string) http.Cookie {
	value := username + ":" + sessID
	return http.Cookie{Name: "sessID", Value: value, MaxAge: 60, HttpOnly: false}
}

func endSessionCookie() http.Cookie {
	return http.Cookie{Name: "sessID", Value: "", MaxAge: 0, HttpOnly: true}
}

func containsForbiddenChar(data []byte) bool {
	forbidden := "[^A-Za-z0-9!@#$%^&*()]"

	m, _ := regexp.Match(forbidden, data)

	return m
}

func readArguments(data []byte) []Arg {
	var args []Arg
	parameters := strings.Split(string(data), ":")
	for _, parameter := range parameters {
		pair := strings.Split(parameter, "=")
		if len(pair) != 2 || containsForbiddenChar([]byte(pair[0])) || containsForbiddenChar([]byte(pair[1])) {
			return nil
		} else {
			var arg Arg
			arg.name = pair[0]
			arg.value = pair[1]
			args = append(args, arg)
		}
	}

	return args
}

func validateSession(db *sql.DB, username string, sessID string) bool {
	// Find all of this user's ongoing sessions
	rows, err := db.Query("SELECT * FROM sessions WHERE username = $1", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Scan current row
		var session Session
		err := rows.Scan(&session.user, &session.id, &session.expires)
		if err != nil {
			if err == sql.ErrNoRows {
				// Return false if user has no ongoing sessions
				return false
			}
		} else {
			// Delete all expired sessions
			if time.Now().After(session.expires) {
				db.Exec("DELETE FROM sessions WHERE sessid = $1", session.id)
			} else if sessID == session.id {
				// Return true if unexpired session matches sessID
				return true
			}
		}
	}

	// If valid session was found return false
	return false
}

func Run(db *sql.DB) {
	/* Logout */
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		// Set sessID to an empty cookie
		cookie := endSessionCookie()
		http.SetCookie(w, &cookie)

		// Retrieve user cookie to end session
		userCookie, err := r.Cookie("sessID")
		if err == nil {
			// If cookie has the right ID extract the data
			if userCookie.Name == "sessID" {
				data := userCookie.Value
				splitData := strings.Split(data, ":")
				if len(splitData) == 2 {
					// Get username and sessID from cookie
					username := splitData[0]
					sessID := splitData[1]
					// Invalidate user session by deleting it from the database
					result, _ := db.Exec("DELETE FROM sessions WHERE username = $1 AND sessid = $2", username, sessID)
					count, _ := result.RowsAffected()
					// If a row was deleted respond with log out confirmation
					if count > 0 {
						w.Write([]byte("Logged out"))
					}
				}
			}
		}

		// If no session was deleted send an empty response that still deletes the cookie
		w.Write([]byte(""))
	})
	/* Secure page */
	http.HandleFunc("/secure", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		// Process the sessID cookie from the request; if it isn't set, just send back an error message
		cookie, err := r.Cookie("sessID")
		if err != nil {
			if err == http.ErrNoCookie {
				w.Write([]byte("Invalid session"))
			} else {
				log.Fatal(err)
			}
		} else {
			cookieValues := strings.Split(cookie.Value, ":")
			if len(cookieValues) == 2 {
				// Obtain username and sessID from cookie and determine if their session is valid
				username, sessID := cookieValues[0], cookieValues[1]
				if validateSession(db, username, sessID) {
					w.Write([]byte("Secure page"))
				} else {
					w.Write([]byte("Invalid session"))
				}
			} else {
				w.Write([]byte("Invalid session"))
			}
		}
	})
	/* Login */
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		// Read request body
		reqBody, err := io.ReadAll(r.Body)
		if err == nil {
			args := readArguments(reqBody)
			if len(args) == 2 &&
				args[0].name == "username" &&
				args[1].name == "password" {
				// Username is set to all lowercase to be case insensitive
				username, password := strings.ToLower(args[0].value), args[1].value

				// Query the database for the specified user
				row := db.QueryRow("SELECT * FROM users WHERE username = $1", username)

				// If no row was found the user does not exist; if it was, compare its hash with the provided password (after hashing it as well)
				var user User
				err := row.Scan(&user.name, &user.hash)
				if err != nil {
					if err == sql.ErrNoRows {
						w.Write([]byte("User not found"))
					} else {
						w.Write([]byte("Error finding user"))
					}
				} else {
					// Hash the provided password to compare to the user data
					match, err := argon2id.ComparePasswordAndHash(password, user.hash)
					if err != nil {
						w.Write([]byte("Error finding user"))
					} else if match {
						// Initialize a new sessions and send back a cookie containing the sessID
						sessID := genID()
						cookie := createCookie(username, sessID)
						http.SetCookie(w, &cookie)
						startSession(db, username, sessID)
						w.Write([]byte("Authenticated"))
					} else {
						w.Write([]byte("Incorrect password"))
					}
				}
			}
		}
	})
	/* Signup */
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		// Read request body
		reqBody, err := io.ReadAll(r.Body)
		if err == nil {
			args := readArguments(reqBody)
			if len(args) == 2 &&
				args[0].name == "username" &&
				args[1].name == "password" {
				// Username is set to all lowercase to be case insensitive
				username, password := strings.ToLower(args[0].value), args[1].value
				// Check if user is already in database
				if UserInDB(db, username) {
					w.Write([]byte("This username is already is use"))
				} else {
					// Add user
					AddUser(db, username, password)
					w.Write([]byte("Account created"))
				}
			}
		}
	})
	http.ListenAndServe(":8000", nil)
}
