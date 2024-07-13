import React, {useEffect, useState} from 'react'
import { useNavigate } from 'react-router-dom'
import useAuth from './../../useAuth'
import './Homepage.css'

export const Homepage = () => {
  const [data, setData] = useState('')
  const navigate = useNavigate()
  const { logout, getData } = useAuth()

  function handleLogout () {
    logout()
    navigate("/login")
  }

  useEffect(() => {
    getData().then((data) => {
      setData(data)
    }, () => {
      logout()
      navigate("/login")
    })
  }, [])
  return (
    <div className='homepage'>
        <div className='bar'>
          <div className='logout' onClick={handleLogout} style={{cursor:'pointer'}}>Log out</div>
          <div className='underline'></div>
        </div>
        <div className='data'>{ data }</div>
    </div>
  )
}

export default Homepage