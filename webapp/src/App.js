import './App.css';
import Nav from './Components/Nav/Nav.jsx'
import { AuthProvider } from './useAuth.js';

function App() {
  return (
    <div>{ AuthProvider( <Nav/> ) }</div>    
  );
}

export default App;
