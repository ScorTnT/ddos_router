import { useState, useEffect, Suspense, lazy } from "react";
import {
  HashRouter,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";
import api from "./api.js";

// React.lazy로 페이지별 분할
const Login = lazy(() => import("./Login.jsx"));
const Dashboard = lazy(() => import("./Dashboard.jsx"));
const APIGuide = lazy(() => import("./APIGuide.jsx"));

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isCheckingAuth, setIsCheckingAuth] = useState(true);

  useEffect(() => {
    // 앱 시작 시 자동 로그인 체크
    const checkAutoLogin = async () => {
      try {
        const isValid = await api.checkAutoLogin();
        if (isValid) {
          console.log('[App] Auto login successful');
          setIsLoggedIn(true);
        } else {
          console.log('[App] No valid session found');
        }
      } catch (error) {
        console.error('[App] Auto login check failed:', error);
      } finally {
        setIsCheckingAuth(false);
      }
    };

    checkAutoLogin();
  }, []);

  // 자동 로그인 체크 중일 때는 로딩 표시
  if (isCheckingAuth) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontSize: '18px' 
      }}>
        로그인 정보 확인 중...
      </div>
    );
  }

  return (
    <HashRouter>
      <Suspense fallback={
        <div style={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          height: '100vh',
          fontSize: '18px' 
        }}>
          페이지 로딩 중...
        </div>
      }>
        <Routes>
          <Route
            index
            element={
              isLoggedIn ? (
                <Navigate to="/dashboard" replace />
              ) : (
                <Navigate to="/login" replace />
              )
            }
          />
          <Route
            path="/api"
            element={
              isLoggedIn ? (
                <APIGuide setIsLoggedIn={setIsLoggedIn} />
              ) : (
                <Navigate to="/login" replace />
              )
            }
          />
          <Route
            path="/login"
            element={
              isLoggedIn ? (
                <Navigate to="/dashboard" replace />
              ) : (
                <Login setIsLoggedIn={setIsLoggedIn} />
              )
            }
          />
          <Route
            path="/dashboard/*"
            element={
              isLoggedIn ? (
                <Dashboard setIsLoggedIn={setIsLoggedIn} />
              ) : (
                <Navigate to="/login" replace />
              )
            }
          />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Suspense>
    </HashRouter>
  );
}

export default App;
