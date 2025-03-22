import {useState, useef} from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Login from './Login.jsx';
import Dashboard from './Dashboard.jsx';
import APIGuide from './APIGuide.jsx';

function App() {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    return (
        <Router>
            <Routes>
                <Route
                    path="/"
                    element={
                        isLoggedIn ? 
                            <Navigate to="/dashboard"/> : 
                            <Navigate to="/login"/>
                    }
                />
                <Route 
                    path="/api" 
                    element={
                        isLoggedIn ? 
                            <APIGuide setIsLoggedIn={setIsLoggedIn} /> :
                            <Navigate to="/login"/> 
                    }
                />
                <Route
                    path="/login"
                    element={
                        isLoggedIn ?
                            <Navigate to="/dashboard"/> :
                            <Login setIsLoggedIn={setIsLoggedIn}/>
                    }
                />
                <Route
                    path="/dashboard/*"
                    element={
                        isLoggedIn ?
                            <Dashboard setIsLoggedIn={setIsLoggedIn}/> :
                            <Navigate to="/login"/>
                    }
                />
            </Routes>
        </Router>
    );
}

export default App;