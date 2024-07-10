import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './../Login/Login.jsx';
import Homepage from './../Homepage/Homepage.jsx'

export const Nav = () => {
    function securePage(Page, Redirect) {
        if (document.cookie.indexOf('sessID=') !== -1) {
            return Page
        } else {
            return Redirect
        }
    }

    return (
        <Router>
            <Routes>
                <Route path="/login" element={ <Login state=""/> }/>
                <Route path="/home" element={ securePage(
                <Homepage/>, 
                <Navigate to="/login" state="You do not have permission to access this page. Please log in."></Navigate>
                ) }/>
                <Route path="*" element={ document.cookie.indexOf('sessID=') === -1 ? <Navigate to="/login" state=""></Navigate> : <Navigate to="/home" state=""></Navigate> }/>
            </Routes>
        </Router>
    )
}

export default Nav