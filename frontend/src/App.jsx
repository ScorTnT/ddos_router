import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Navigate, Routes } from 'react-router-dom';
import Login from './Login';
import Config from './Config';
import NetworkSetting from './NetworkSetting';
function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route
            path="/"
            element={isLoggedIn ? <NetworkSetting/> : <Login/>}
          />
          <Route
            path="/login"
            element={
              isLoggedIn ?
              <Config /> :
              <Login setIsLoggedIn={setIsLoggedIn} />
            }
          />
          <Route
            path="/config"
            element={isLoggedIn ? <Config/> : <Login />}
          />
        </Routes>
      </div>
    </Router>
  );
}
export default App;