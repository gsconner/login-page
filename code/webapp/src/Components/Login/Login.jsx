import React, {useState} from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import useAuth from './../../useAuth'
import './Login.css'

export const Login = () => {
    const [username, setName] = useState('')
    const [password, setPass] = useState('')
    const [message, setMsg] = useState(useLocation().state)
    const [loading, setLoading] = useState(false)
    const navigate = useNavigate();
    const { login } = useAuth()

    const handleSubmit = (event) => {
        event.preventDefault();
        
        const forbid = /[\<\>\,\.\?\[\]\|\{\}\:\=\;\'\"\/\\\-\+]/
        if (username.match(forbid) || password.match(forbid)) {
            setMsg("You may not create a username or password with the following characters: < > , . ? [ ] | { } : = ; ' \" / \\ - +")
        } else {
            setLoading(true)
            login(username, password).then(() => {
                navigate("/home")
            }, (data) => {
                setMsg(data)
            }).finally(() => {
                setPass('')
                setLoading(false)
            })
        }
    }
    return (
        <div className='loginpage'>
            <div className='message'>{message}</div>
            <div className='window'>
                <div className='header'>Login</div>
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
                    </div>
                    <div className='submit'>
                        <button type='submit' className='button'>Submit</button>
                        <div className='loading'>{loading===true?"Loading...":""}</div>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default Login