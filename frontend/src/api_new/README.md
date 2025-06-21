# Router API Client Library

이 라이브러리는 DDoS 라우터 API와 상호작용하기 위한 JavaScript 클라이언트입니다.

## 설치 및 사용법

### 기본 사용법

```javascript
import apiClient from './api_new/index.js';

// 베이스 URL 설정 (선택사항)
apiClient.setBaseURL('http://localhost:8080');

// 로그인
try {
  const loginResponse = await apiClient.login('username', 'password');
  console.log('로그인 성공:', loginResponse);
} catch (error) {
  console.error('로그인 실패:', error);
}

// 라우터 정보 가져오기
try {
  const routerInfo = await apiClient.information.getRouterInfo();
  console.log('라우터 정보:', routerInfo);
} catch (error) {
  console.error('라우터 정보 가져오기 실패:', error);
}
```

### 인증 (Authentication)

```javascript
// 로그인
await apiClient.login('admin', 'password');

// 로그아웃
await apiClient.logout();

// 인증 상태 확인
const isLoggedIn = apiClient.isAuthenticated();
```

### 정보 조회 (Information)

```javascript
// 라우터 정보
const routerInfo = await apiClient.information.getRouterInfo();

// 네이버 정보 (ARP 테이블)
const neighbors = await apiClient.information.getNeighbors();

// 연결 정보
const connections = await apiClient.information.getConnections();
```

### 설정 조회 (Configuration)

```javascript
// LAN 설정
const lanConfig = await apiClient.config.getLANConfig();

// WAN 설정
const wanConfig = await apiClient.config.getWANConfig();

// 모든 설정
const allConfigs = await apiClient.config.getAllConfigs();
```

### 보안 관리 (Protection)

```javascript
// 보호 목록 조회
const protectionList = await apiClient.protection.getProtectionList();

// IP 차단
await apiClient.protection.blockIP('192.168.1.100');

// IP 차단 해제
await apiClient.protection.unblockIP('192.168.1.100');

// IP 차단 상태 토글
await apiClient.protection.toggleIPBlock('192.168.1.100', true); // 차단
await apiClient.protection.toggleIPBlock('192.168.1.100', false); // 해제
```

### 에러 핸들링

```javascript
import { formatErrorMessage, isAuthError } from './api_new/utils.js';

try {
  await apiClient.information.getRouterInfo();
} catch (error) {
  const errorMessage = formatErrorMessage(error);
  
  if (isAuthError(error)) {
    console.log('인증 오류:', errorMessage);
    // 로그인 페이지로 리다이렉트
  } else {
    console.log('일반 오류:', errorMessage);
  }
}
```

### 유틸리티 함수

```javascript
import { isValidIP, retryWithBackoff, debounce } from './api_new/utils.js';

// IP 주소 검증
const isValid = isValidIP('192.168.1.1'); // true

// 재시도 기능
const result = await retryWithBackoff(
  () => apiClient.information.getRouterInfo(),
  3, // 최대 3번 재시도
  1000 // 1초 기본 지연
);

// 디바운스 함수
const debouncedSearch = debounce((query) => {
  // 검색 로직
}, 300);
```

## API 엔드포인트

### 인증
- `POST /api/auth/login` - 로그인
- `POST /api/auth/logout` - 로그아웃

### 정보
- `GET /api/information` - 라우터 정보
- `GET /api/information/neighbors` - ARP 테이블
- `GET /api/information/connections` - 연결 정보

### 설정
- `GET /api/config/lan` - LAN 설정
- `GET /api/config/wan` - WAN 설정

### 보안
- `GET /api/protection` - 보호 목록
- `POST /api/protection/ip/block` - IP 차단
- `POST /api/protection/ip/unblock` - IP 차단 해제

## 파일 구조

```
api_new/
├── index.js          # 메인 API 클라이언트
├── base.js           # 베이스 API 클래스
├── auth.js           # 인증 API
├── information.js    # 정보 조회 API
├── config.js         # 설정 API
├── protection.js     # 보안 API
├── utils.js          # 유틸리티 함수
└── README.md         # 사용 설명서
```

## 주의사항

1. 모든 API 요청은 비동기적으로 처리됩니다.
2. 인증이 필요한 API는 먼저 로그인을 해야 합니다.
3. 세션 ID는 자동으로 관리되며, 로그인 후 모든 요청에 포함됩니다.
4. 에러 처리를 위해 try-catch 블록을 사용하세요.
