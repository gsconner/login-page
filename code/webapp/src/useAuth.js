import React, { useState } from 'react'

const authContext = React.createContext()

function useAuth() {
    const [authed, setAuthed] = useState(false)

    return {
        authed,
        login(username, password) {
            return new Promise((resolve, reject) => {
                const reqBody = 'username=' + username + ':password=' + password
                
                fetch('http://mylogin:8000/login', {
                    method: 'POST',
                    credentials: 'include',
                    body: reqBody
                }).then(response => {
                    if (response.ok) {
                        response.text().then(msg => {
                            if (msg === "Authenticated") {
                                setAuthed(true)
                                resolve()
                            } else {
                                reject(msg)
                            }
                        })
                    } else {
                        reject("Bad response from server")
                    }
                }).catch(error => {
                    console.log(error)
                    reject("No response from server")
                })
            })
        },
        logout() {
            fetch('http://mylogin:8000/logout', {
                method: 'POST',
                credentials: 'include'
            }).catch(error => console.log(error))
            setAuthed(false)
        },
        getData() {
            return new Promise((resolve, reject) => {
                fetch('http://mylogin:8000/secure', {
                    method: 'GET',
                    credentials: 'include'
                }).then(response => {
                    if (response.ok) {
                        response.text().then(data => {
                            if (data !== '') {
                                resolve(data)
                            } else {
                                reject(data)
                            }
                        })
                    }
                }).catch(error => {
                    console.log(error)
                })
            })
        }
    }
}

export function AuthProvider(children) {
    const auth = useAuth()

    return <authContext.Provider value={auth}>{children}</authContext.Provider>
}

export default function AuthConsumer() {
    return React.useContext(authContext)
}