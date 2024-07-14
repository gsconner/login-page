package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
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
	ip      string
	expires time.Time
}

type Lock struct {
	failedAttempts int
	locked         bool
	expires        time.Time
}

var lockouts map[string](map[string]*Lock)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func initLock() *Lock {
	lock := new(Lock)
	lock.failedAttempts = 0
	lock.locked = false
	lock.expires = time.Now().UTC().Add(time.Minute)

	return lock
}

func startSession(db *sql.DB, username string, sessID string, clientIP string) {
	_, err := db.Exec("INSERT INTO sessions (username, sessid, ip, expires) VALUES ($1, $2, $3, $4)", username, sessID, clientIP, time.Now().UTC().Add(time.Minute).Format(time.DateTime))
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

func validateSession(db *sql.DB, username string, sessID string, clientIP string) bool {
	// Find all of this user's ongoing sessions
	rows, err := db.Query("SELECT * FROM sessions WHERE username = $1", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		// Scan current row
		var session Session
		err := rows.Scan(&session.user, &session.id, &session.ip, &session.expires)
		if err != nil {
			if err == sql.ErrNoRows {
				// Return false if user has no ongoing sessions
				return false
			}
		} else {
			// Delete all expired sessions
			if time.Now().After(session.expires) {
				db.Exec("DELETE FROM sessions WHERE sessid = $1", session.id)
			} else if clientIP == session.ip && sessID == session.id {
				// If session is not expired return true
				return true
			}
		}
	}

	// If valid session was found return false
	return false
}

func getClientAddr(r *http.Request) string {
	clientAddr := r.Header.Get("X-FORWARDED-FOR")
	if clientAddr == "" {
		return r.RemoteAddr
	} else {
		return clientAddr
	}
}

func isLockedOut(clientIP string, username string) (bool, *time.Time) {
	locks, ok := lockouts[clientIP]
	if ok {
		lock, ok := locks[username]
		if ok {
			if lock.locked == true {
				if time.Now().After(lock.expires) {
					locks[username] = initLock()
				} else {
					return true, &lock.expires
				}
			}
		}
	}

	return false, nil
}

func recordInvalidLogin(clientIP string, username string) (bool, int) {
	locks, ok := lockouts[clientIP]
	if ok {
		lock, ok := locks[username]
		if !ok {
			locks[username] = initLock()
			lock = locks[username]
			lock.failedAttempts += 1
			return false, lock.failedAttempts
		} else {
			if time.Now().After(lock.expires) {
				locks[username] = initLock()
				lock = locks[username]
			} else {
				lock.expires = time.Now().UTC().Add(time.Minute)
			}

			lock.failedAttempts += 1
			if lock.failedAttempts < 10 {
				return false, lock.failedAttempts
			} else {
				lock.locked = true
				return true, 0
			}
		}
	} else {
		lockouts[clientIP] = make(map[string]*Lock)
		locks := lockouts[clientIP]
		locks[username] = initLock()
		lock := locks[username]
		lock.failedAttempts += 1

		return false, lock.failedAttempts
	}
}

func Run(db *sql.DB) {
	// Initialize server variables
	lockouts = make(map[string](map[string]*Lock))

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
			// Cient IP and lockouts have been disabled for now
			//clientAddr := getClientAddr(r)
			//fmt.Printf("%s\n", clientAddr)
			clientAddr := "192.0.0.1:0000"
			addrSplit := strings.Split(clientAddr, ":")
			if len(addrSplit) == 2 {
				clientIP := addrSplit[0]
				cookieValues := strings.Split(cookie.Value, ":")
				if len(cookieValues) == 2 {
					// Obtain username and sessID from cookie and determine if their session is valid
					username, sessID := cookieValues[0], cookieValues[1]
					if validateSession(db, username, sessID, clientIP) {
						w.Write([]byte("Secure page"))
					} else {
						w.Write([]byte("Invalid session"))
					}
				} else {
					w.Write([]byte("Invalid session"))
				}
			}
		}
	})
	/* Login */
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		// Read request body
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			reqSplit := strings.Split(string(reqBody), ":")
			if len(reqSplit) == 2 {
				// Split request data into username and password
				args := [2][]string{strings.Split(reqSplit[0], "="), strings.Split(reqSplit[1], "=")}
				if len(args[0]) == 2 && args[0][0] == "username" &&
					len(args[1]) == 2 && args[1][0] == "password" {
					// Username is set to all lowercase to be case insensitive
					username, password := strings.ToLower(args[0][1]), args[1][1]

					// Cient IP and lockouts have been disabled for now
					//clientAddr := getClientAddr(r)
					clientAddr := "192.0.0.1:0000"
					addrSplit := strings.Split(clientAddr, ":")
					if len(addrSplit) == 2 {
						clientIP := addrSplit[0]
						// Make sure client is not locked out from attempting login on this account (no longer used)
						/*lock, expires := isLockedOut(clientIP, username)
						if lock {
							if expires == nil {
								w.Write([]byte("Error receiving client ip addr"))
							}
							w.Write([]byte(fmt.Sprintf("You are locked from logging into this account until %s", expires.Format(time.Stamp))))
						} else*/{
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
									startSession(db, username, sessID, clientIP)
									w.Write([]byte("Authenticated"))
								} else {
									// Record invalid login attempt and lock the user out if too many attempts (not used anymore)
									/*locked, attempts := recordInvalidLogin(clientIP, username)
									if locked {
										w.Write([]byte("You have been locked for too many failed attempts"))
									} else*/{
										//w.Write([]byte(fmt.Sprintf("Incorrect password. You have %d attempts remaining.", (10 - attempts))))
										w.Write([]byte("Incorrect password"))
									}
								}
							}
						}
					}
				}
			}
		}
	})
	http.ListenAndServe(":8000", nil)
}
