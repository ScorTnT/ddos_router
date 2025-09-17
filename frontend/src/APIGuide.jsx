import {
  Box,
  Typography,
  Stack,
  Paper,
  List,
  ListItemButton,
  ListItemText,
} from "@mui/material";
import { useState } from "react";

// 2.0 버전의 API 가이드 내용만 남겼습니다.
const GuideContent = {
  "시작하기": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          본 문서는 Ddos 라우터 API 개발 연동을 위한 공식 가이드 문서입니다. <br />
          API는 **Axios 기반**의 단일 클라이언트 파일(`api.js`)로 통합되어, 효율적인 관리를 제공합니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · API 규격
        </Typography>
        <Typography>
          **RESTful** 방식으로 호출이 가능하며 **GET/POST**만 지원합니다. <br />
          인증이 필요한 API는 먼저 로그인을 해야 합니다. 로그인 성공 시 발급되는 세션 ID를 `X-Session-ID` 헤더로 자동 전송합니다. <br />
          응답 규약은 `{` status: "success", data: ... `}` 또는 `{` status: "error", message: "..." `}` 입니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 호출 예제
        </Typography>
        <Typography>
          아래는 `api.js` 파일을 사용한 로그인 예제입니다.
        </Typography>
        <pre>
          <code>
            {`import api from './api.js';

async function performLogin() {
    try {
        const username = 'root';
        const password = 'password';
        const loginSuccessful = await api.login(username, password, { remember: true });
        if (loginSuccessful) {
            console.log('로그인 성공!');
            const info = await api.getInformation();
            console.log('라우터 정보:', info);
        } else {
            console.log('로그인 실패!');
        }
    } catch (error) {
        console.error('로그인 중 오류 발생:', error);
    }
}
performLogin();`}
          </code>
        </pre>
      </Paper>
    </Stack>
  ),

  "인증": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          Ddos 라우터 계정에 로그인, 로그아웃, 자동 로그인 상태를 확인하는 API입니다.
          로그인 성공 시 발급되는 세션 ID는 `X-Session-ID` 헤더에 자동으로 포함됩니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · API 엔드 포인트
        </Typography>
        <Typography>
          **POST** `/api/auth/login` - 로그인 <br />
          **POST** `/api/auth/logout` - 로그아웃 <br />
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · Request Parameters (`/api/auth/login`)
        </Typography>
        <pre>
          <code>
            {`{
    "username": "string",
    "password": "string"
}`}
          </code>
        </pre>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · Received Message (`/api/auth/login`)
        </Typography>
        <pre>
          <code>
            {`{
    "status": "success",
    "data": {
        "message": "로그인 성공",
        "session_id": "string",
        "expires_at": "string"
    }
}`}
          </code>
        </pre>
      </Paper>
    </Stack>
  ),

  "정보 조회": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          라우터의 현재 상태 및 연결된 장치에 대한 정보를 조회합니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · API 엔드 포인트
        </Typography>
        <Typography>
          **GET** `/api/information` - 라우터 전체 정보 <br />
          **GET** `/api/information/neighbors` - 이웃(Neighbors) 장치 정보 <br />
          **GET** `/api/information/connections` - 현재 연결 정보 <br />
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · Response Sample Code (Json)
        </Typography>
        <Typography>
          **`GET` `/api/information`**
        </Typography>
        <pre>
          <code>
            {`{
    "status": "success",
    "data": {
        "Hostname": "DDOS-Router",
        "OS Version": "1.0.0",
        "MAC Address": "00:11:22:33:44:55"
    }
}`}
          </code>
        </pre>
        <Typography sx={{ mt: 2 }}>
          **`GET` `/api/information/neighbors`**
        </Typography>
        <pre>
          <code>
            {`{
    "status": "success",
    "data": [
        { "ip": "192.168.1.10", "mac": "AA:BB:CC:DD:EE:FF" }
    ]
}`}
          </code>
        </pre>
      </Paper>
    </Stack>
  ),

  "설정 조회": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          라우터의 WAN 및 LAN 네트워크 설정을 조회하고 변경합니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · API 엔드 포인트
        </Typography>
        <Typography>
          **GET** `/api/config/wan` - WAN 설정 조회 <br />
          **POST** `/api/config/wan` - WAN 설정 업데이트 <br />
          **GET** `/api/config/lan` - LAN 설정 조회 <br />
          **POST** `/api/config/lan` - LAN 설정 업데이트 <br />
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · Request Parameters (`/api/config/wan`)
        </Typography>
        <pre>
          <code>
            {`{
    "connection_type": "dhcp" | "static",
    "ip_addr": "string",
    "netmask": "string",
    "gateway": "string",
    "dns_list": ["string", "string"],
    "mac_addr": "string",
    "mtu": "number",
    "is_custom_dns": "boolean",
    "is_custom_mac": "boolean",
    "manualMtu": "boolean"
}`}
          </code>
        </pre>
      </Paper>
    </Stack>
  ),

  "보안 관리": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          악성 IP를 차단하거나 차단된 IP를 허용하는 API입니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · API 엔드 포인트
        </Typography>
        <Typography>
          **GET** `/api/protection` - 차단된 IP 목록 조회 <br />
          **POST** `/api/protection/ip/block` - IP 차단 (쿼리 파라미터로 IP 전달) <br />
          **POST** `/api/protection/ip/unblock` - IP 차단 해제 (쿼리 파라미터로 IP 전달) <br />
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · Request Parameters (`/api/protection/ip/block`)
        </Typography>
        <Typography>
          - **Query String**: `?ip=string&is_permanent=boolean`<br />
          - **ip**: 차단할 IP 주소 (필수)<br />
          - **is_permanent**: 영구 차단 여부 (선택, 기본값: false)
        </Typography>
      </Paper>
    </Stack>
  ),

  "에러 헨들링": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          `api.js` 파일의 응답 인터셉터를 통해 모든 API 요청에 대한 에러를 중앙에서 처리합니다.
          특히, 세션이 만료된 경우(HTTP 401) 자동으로 로그아웃을 처리합니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 에러 코드
        </Typography>
        <Typography>
          - **401 Unauthorized**: 세션 ID가 유효하지 않거나 만료된 경우.<br />
          - **404 Not Found**: 요청한 엔드포인트가 존재하지 않거나, (특정 API의 경우) 리소스가 없는 경우.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · API 응답 규약
        </Typography>
        <Typography>
          오류 응답은 다음 형식을 따릅니다.
        </Typography>
        <pre>
          <code>
            {`{
    "status": "error",
    "message": "오류 메시지"
}`}
          </code>
        </pre>
      </Paper>
    </Stack>
  ),

  "유틸리티": (
    <Stack spacing={3}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 개요
        </Typography>
        <Typography>
          API 클라이언트 내에서 세션 관리 및 기타 보조 기능을 제공합니다.
        </Typography>
      </Paper>
      <Paper sx={{ p: 3 }}>
        <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
          · 주요 기능
        </Typography>
        <Typography>
          - `setSession(id, persist)`: 세션 ID를 `sessionStorage` 또는 `localStorage`에 저장합니다. <br />
          - `clearSession()`: 저장된 세션 ID를 삭제합니다. <br />
          - `getSessionId()`: 저장된 세션 ID를 가져옵니다.
        </Typography>
      </Paper>
    </Stack>
  ),
};

const ApiGuide = () => {
  const [selectedMenu, setSelectedMenu] = useState("시작하기");
  const menuItems = Object.keys(GuideContent);
  const currentContent = GuideContent[selectedMenu];

  return (
    <Box display="flex" sx={{ p: 3 }}>
      <Box
        sx={{
          minWidth: 250,
          borderRight: "1px solid #ddd",
          pr: 2,
        }}
      >
        <Typography variant="h6" gutterBottom>
          API 개발가이드
        </Typography>
        <List dense>
          {menuItems.map((item) => (
            <ListItemButton
              key={item}
              selected={selectedMenu === item}
              onClick={() => setSelectedMenu(item)}
            >
              <ListItemText
                primary={item}
                primaryTypographyProps={{
                  sx: {
                    color: selectedMenu === item ? "primary.main" : "inherit",
                  },
                }}
              />
            </ListItemButton>
          ))}
        </List>
      </Box>

      <Box sx={{ flexGrow: 1, pl: 4 }}>
        <Typography variant="h5" gutterBottom>
          {selectedMenu}
        </Typography>
        {currentContent}
      </Box>
    </Box>
  );
};

export default ApiGuide;