import React, {useState} from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import useAuth from './../../useAuth'
import './Login.css'

export const Login = () => {
    const [username, setName] = useState('')
    const [password, setPass] = useState('')
    const [confirmPass, setConfirm] = useState('')
    const [message, setMsg] = useState(useLocation().state)
    const [loading, setLoading] = useState(false)
    const [mode, setMode] = useState(false)
    const navigate = useNavigate();
    const { login, signup } = useAuth()

    const handleSubmit = (event) => {
        event.preventDefault();
        
        const forbid = /[^A-Za-z0-9!@#$%^&*()]/
        if (username.match(forbid) || password.match(forbid)) {
            setMsg("Invalid character")
        } else {
            if (mode) {
                if (password === confirmPass) {
                    setLoading(true)
                    setMsg('')
                    signup(username, password, confirmPass).then((data) => {
                        setMsg(data)
                    }, (data) => {
                        setMsg(data)
                    }).finally(() => {
                        setName('')
                        setPass('')
                        setConfirm('')
                        setLoading(false)
                    })
                } else {
                    setMsg("Password does not match")
                }
            } else {
                setLoading(true)
                setMsg('')
                login(username, password).then(() => {
                    navigate("/home")
                }, (data) => {
                    setMsg(data)
                }).finally(() => {
                    setName('')
                    setPass('')
                    setConfirm('')
                    setLoading(false)
                })
            }
        }
    }
    return (
        <div className='loginpage'>
            <div className='message'>{message}</div>
            <div className='window'>
                <div className='header'>{mode ? "Sign Up" : "Login" }</div>
                <form onSubmit={handleSubmit}>
                    <div className='inputs'>
                        <div className='input'>
                            <input 
                                type='username' 
                                required
                                placeholder='Username' 
                                value={username} 
                                onChange={(e) => setName(e.target.value)}
                            ></input>
                        </div>
                        <div className='input'>
                            <input 
                                type='password'
                                required
                                placeholder='Password'
                                value={password} 
                                onChange={(e) => setPass(e.target.value)}
                            ></input>
                        </div>
                        {mode ?
                            <div className='input'>
                                <input 
                                    type='password'
                                    required
                                    placeholder='Confirm Password'
                                    value={confirmPass} 
                                    onChange={(e) => setConfirm(e.target.value)}
                                ></input>
                            </div>
                            :null
                        }
                    </div>
                    <div className='submit'>
                        <button type='submit' className='button'>Submit</button>
                        <div className='loading'>{loading===true?"Loading...":""}</div>
                    </div>
                </form>
            </div>
            <div className='modechange'>
                <h1 onClick={() => setMode(!mode)}>
                    <div className='signup'>{mode ? "Login" : "Sign Up" }</div>
                </h1>
            </div>
        </div>
    )
}

export default Login