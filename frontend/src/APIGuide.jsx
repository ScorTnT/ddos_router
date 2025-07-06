import {
  Box,
  Typography,
  Stack,
  Paper,
  MenuItem,
  Select,
  List,
  ListItemButton,
  ListItemText,
  Divider,
} from "@mui/material";
import { useState } from "react";
// 1) 버전·메뉴별 JSX를 한곳에 모아둡니다.
const MenuMap = {
  "1.0": {
    "시작하기" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 개요 (v1.0)
          </Typography>
          <Typography>
            본 문서는 Ddos 라우터 API 개발 연동을 위한 가이드 문서입니다. <br />
            이는 개발 초기의 구버전 API <br />
          </Typography>
        </Paper>
        {/* v1.0만의 추가 설명 */}
      </Stack>
    ),
    // 1.0의 다른 메뉴
  },

  "2.0": {
    "시작하기" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 개요 (v2.0)
          </Typography>
          <Typography>
            본 문서는 Ddos 라우터 API 개발 연동을 위한 가이드 문서입니다. <br />

          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 규격
          </Typography>
          <Typography>
            RESTful 방식으로 호출이 가능하며 Method는 GET/POST만 지원합니다. <br />
            JSON과 XML 규격을 제공합니다. <br />
            인증이 필요한 API는 먼저 로그인을 해야 합니다. <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 주의 사항
          </Typography>
          <Typography>
            모든 API 요청은 비동기적으로 처리됩니다. <br />
            세션 ID는 자동으로 관리되며, 로그인 후 모든 요청에 포함됩니다. <br />
            에러 처리를 위해 try - catch 블록을 사용하세요. <br /> 
          </Typography>
        </Paper>
        
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 호출 예제
          </Typography>
          <Typography>
            호출 예제 코드 <br /> 
          </Typography>
        </Paper>              
      </Stack>
    ),

    "인증" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 인증
          </Typography>
          <Typography>
            Ddos 라우터 계정에 로그인 / 로그아웃 하는 API 입니다.
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 엔드 포인트
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Parameters
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Received Message
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Response Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>
      </Stack>
    ),

    "정보 조회" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 인증
          </Typography>
          <Typography>
            Ddos 라우터 계정에 로그인 / 로그아웃 하는 API 입니다.
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 엔드 포인트
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Parameters
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Received Message
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Response Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>
      </Stack>
    ),

    "설정 조회" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 인증
          </Typography>
          <Typography>
            Ddos 라우터 계정에 로그인 / 로그아웃 하는 API 입니다.
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 엔드 포인트
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Parameters
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Received Message
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Response Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>
      </Stack>
    ),

    "보안 관리" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 인증
          </Typography>
          <Typography>
            Ddos 라우터 계정에 로그인 / 로그아웃 하는 API 입니다.
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 엔드 포인트
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Parameters
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Received Message
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Response Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>
      </Stack>
    ),

    "에러 헨들링" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 인증
          </Typography>
          <Typography>
            Ddos 라우터 계정에 로그인 / 로그아웃 하는 API 입니다.
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 엔드 포인트
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Parameters
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Received Message
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Response Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>
      </Stack>
    ),

    "유틸리티" : (
      <Stack spacing={3}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · 인증
          </Typography>
          <Typography>
            Ddos 라우터 계정에 로그인 / 로그아웃 하는 API 입니다.
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · API 엔드 포인트
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Parameters
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Received Message
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Request Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>

        <Paper sx={{ p: 3 }}>
          <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
            · Response Sample Code (Json)
          </Typography>
          <Typography>
            POST /api/auth/login - 로그인 <br />
            POST /api/auth/logout - 로그아웃 <br />
          </Typography>
        </Paper>
      </Stack>
    ),                   
  },
};

const ApiGuide = () => {
  const [version, setVersion] = useState("2.0");
  const [selectedMenu, setSelectedMenu] = useState("시작하기");
  const menuItems = Object.keys(MenuMap[version]);
  const currentContent = MenuMap[version][selectedMenu];
  const handleVersion = (e) => {
    const newVersion = e.target.value;
    setVersion(newVersion);
    setSelectedMenu(Object.keys(MenuMap[newVersion])[0]);
  }

  return (
    <Box display="flex" sx={{ p: 3 }}>
      {/* 왼쪽 메뉴 영역 */}
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

        <Select
          fullWidth
          size="small"
          value={version}
          onChange={handleVersion}
          sx={{ mb: 2 }}
        >
          <MenuItem value="1.0">1.0</MenuItem>
          <MenuItem value="2.0">2.0</MenuItem>
        </Select>

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

      {/* 오른쪽 본문 영역 */}
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