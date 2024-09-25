import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Navigate, Routes } from 'react-router-dom';
import Login from './Login';
import Config from './Config';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route
            path="/"
            element={isLoggedIn ? <Navigate to="/config" /> : <Navigate to="/login" />}
          />
          <Route
            path="/login"
            element={
              isLoggedIn ?
              <Navigate to="/config" /> :
              <Login setIsLoggedIn={setIsLoggedIn} />
            }
          />
          <Route
            path="/config"
            element={isLoggedIn ? <Config /> : <Navigate to="/login" />}
          />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
