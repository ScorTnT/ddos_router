import React, { useState } from 'react';
import { attemptLogin } from './checkAccount.js';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginResult, setLoginResult] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const result = await attemptLogin(username, password);
    setLoginResult(result);
  };

  return (
    <div className="App">
      <h1>Login Test</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">Username:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit">Login</button>
      </form>
      {loginResult !== null && (
        <p>
          Login result: {loginResult ? 'Success' : 'Failure'}
        </p>
      )}
    </div>
  );
}

export default Login;
