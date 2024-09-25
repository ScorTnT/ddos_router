import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Link, Routes } from 'react-router-dom';
import Login from './Login';

function Home() {
  return <h1>Welcome to the Home Page</h1>;
}

function App() {
  let [routing, setRouting] = useState(['networkSetting', 'userSetting', 'intranetSetting']);
  return (
    <Router>
      <div className="App">

        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
          </ul>
        </nav>

        <Routes>
          <Route exact path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
        </Routes>
      </div>

    </Router>
  );
}

export default App;
