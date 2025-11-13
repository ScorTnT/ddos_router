// Simple axios-based API client for DDOS Router backend
// - 환경변수: `VITE_API_URL` 를 baseURL로 사용합니다.
// - 인증: 로그인 성공 시 발급되는 세션 ID를 `X-Session-ID` 헤더로 전송합니다.
// - 응답 규약: { status: "success", data: ... } / { status: "error", message: "..." }

import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_API_URL;
const SESSION_KEY = "session_id";

// 디버깅: 환경변수 확인
console.log('[API] VITE_API_URL:', API_BASE_URL);
console.log('[API] Resolved baseURL:', API_BASE_URL || undefined);

// axios 인스턴스
const http = axios.create({
  baseURL: API_BASE_URL || undefined, // 빈 문자열이면 undefined로 (상대 경로 사용)
  headers: { "Content-Type": "application/json" },
  timeout: 15000,
});

// 요청 인터셉터: 세션 헤더 자동 첨부
http.interceptors.request.use((config) => {
  const sidSession = sessionStorage.getItem(SESSION_KEY);
  const sidLocal = localStorage.getItem(SESSION_KEY);
  const sid = sidSession || sidLocal;
  
  console.log('[API] Request interceptor - SessionStorage:', sidSession);
  console.log('[API] Request interceptor - LocalStorage:', sidLocal);
  console.log('[API] Request interceptor - Final sid:', sid);
  
  if (sid) {
    config.headers["X-Session-ID"] = sid;
    console.log('[API] Adding session header:', sid);
    console.log('[API] Final headers:', config.headers);
  } else {
    console.warn('[API] No session ID found for request to:', config.url);
  }
  
  console.log('[API] Request:', config.method?.toUpperCase(), config.url);
  return config;
});

// 응답 인터셉터: 에러 로깅
http.interceptors.response.use(
  (response) => {
    console.log('[API] Response:', response.status, response.config.url);
    return response;
  },
  (error) => {
    console.error('[API] Error:', error.response?.status, error.config?.url, error.response?.data);
    return Promise.reject(error);
  }
);

// 공통 응답 언래핑
function unwrap(response) {
  const payload = response?.data;
  if (payload && payload.status === "success") return payload.data;
  const message = payload?.message || "요청 처리 중 오류가 발생했습니다.";
  const error = new Error(message);
  error.response = response;
  throw error;
}

// 헬스체크 (루트 엔드포인트)
async function ping() {
  const res = await http.get("/");
  return unwrap(res); // { message: "DDOS Router API" }
}

async function isOnline() {
  try {
    await ping();
    return true;
  } catch {
    return false;
  }
}

// Auth
async function login(username, password, options = { remember: true }) {
  const res = await http.post("/api/auth/login", { username, password });
  const data = unwrap(res); // { message, session_id, expires_at }
  console.log('[API] Login response data:', data);
  if (data?.session_id) {
    console.log('[API] Setting session:', data.session_id);
    setSession(data.session_id, !!options.remember);
  } else {
    console.error('[API] No session_id in response:', data);
  }
  return data;
}

// 자동 로그인 체크 (앱 시작 시 호출)
async function checkAutoLogin() {
  const sessionId = getSessionId();
  if (!sessionId) {
    console.log('[API] No stored session found');
    return false;
  }
  
  console.log('[API] Found stored session, checking validity:', sessionId);
  try {
    // 세션이 유효한지 확인 (간단한 API 호출로 테스트)
    await getInformation();
    console.log('[API] Stored session is valid');
    return true;
  } catch (error) {
    console.log('[API] Stored session is invalid, clearing:', error.message);
    clearSession();
    return false;
  }
}

async function logout() {
  try {
    await http.post("/api/auth/logout");
  } catch (error) {
    console.error('[API] Logout request failed:', error);
  } finally {
    clearSession();
  }
}

function getSessionId() {
  return sessionStorage.getItem(SESSION_KEY) || localStorage.getItem(SESSION_KEY) || "";
}

// 세션 관리 유틸
function setSession(id, persist = false) {
  console.log('[API] setSession called:', id, 'persist:', persist);
  sessionStorage.setItem(SESSION_KEY, id || "");
  if (persist) localStorage.setItem(SESSION_KEY, id || "");
  else localStorage.removeItem(SESSION_KEY);
  console.log('[API] Session stored. SessionStorage:', sessionStorage.getItem(SESSION_KEY));
}

function clearSession() {
  sessionStorage.removeItem(SESSION_KEY);
  localStorage.removeItem(SESSION_KEY);
}

// Information
async function getInformation() {
  const res = await http.get("/api/information");
  return unwrap(res);
}

async function getNeighbors() {
  const res = await http.get("/api/information/neighbors");
  return unwrap(res);
}

async function getConnections() {
  const res = await http.get("/api/information/connections");
  return unwrap(res);
}

// Config: WAN
async function getWANConfig() {
  const res = await http.get("/api/config/wan");
  return unwrap(res);
}

async function updateWANConfig(data) {
  const res = await http.post("/api/config/wan", data);
  return unwrap(res);
}

// Config: LAN
async function getLANConfig() {
  const res = await http.get("/api/config/lan");
  return unwrap(res);
}

async function updateLANConfig(data) {
  const res = await http.post("/api/config/lan", data);
  return unwrap(res);
}

// Protection
async function getProtection() {
  try {
    const res = await http.get("/api/protection");
    return unwrap(res);
  } catch (err) {
    // 서버가 빈 목록일 때 404를 돌려줌 → 빈 배열로 치환
    const status = err?.response?.status;
    if (status === 404) return [];
    throw err;
  }
}

async function blockIP(ip, isPermanent = false) {
  if (!ip) throw new Error("IP가 필요합니다.");
  // 백엔드 스펙: 쿼리스트링 ?ip=...
  const res = await http.post("/api/protection/ip/block", null, { params: { ip, is_permanent: isPermanent } });
  return unwrap(res);
}

async function unblockIP(ip) {
  if (!ip) throw new Error("IP가 필요합니다.");
  const res = await http.post("/api/protection/ip/unblock", null, { params: { ip } });
  return unwrap(res);
}

async function getWhiteList() {
  try{
    const res = await http.get("/api/protection/whitelist");
    return unwrap(res);
  } catch (err) {
    const status = err?.response?.status;
    if(status ===404) return [];
    throw err;
  }
}

async function addWhiteList(ip) {
  if (!ip) throw new Error("IP가 필요합니다.");
  const res = await http.post("/api/protection/whitelist/add", null, { params: { ip } });
  return unwrap(res);
}

async function removeWhiteList(ip) {
  if (!ip) throw new Error("IP가 필요합니다.");
  const res = await http.post("/api/protection/whitelist/remove", null, { params: { ip } });
  return unwrap(res);
}

// 사용 예)
// import api from "./api";
// await api.login("root", "password", { remember: true });
// const info = await api.getInformation();

const api = {
  // health
  ping,
  isOnline,
  // auth
  login,
  logout,
  checkAutoLogin,
  getSessionId,
  // information
  getInformation,
  getNeighbors,
  getConnections,
  // config
  getWANConfig,
  updateWANConfig,
  getLANConfig,
  updateLANConfig,
  // protection
  getProtection,
  blockIP,
  unblockIP,
  // whiteList
  getWhiteList,
  addWhiteList,
  removeWhiteList,
  // low-level
  http,
  setSession,
  clearSession,
};

export default api;
export { ping, isOnline };