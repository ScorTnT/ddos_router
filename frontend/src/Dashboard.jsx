import { useEffect, useState } from 'react';
import {
    AppBar,
    Box,
    Button,
    Card,
    CardContent,
    FormControl,
    FormControlLabel,
    IconButton,
    InputLabel,
    MenuItem,
    Paper,
    Select,
    Stack,
    Switch,
    Tab,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Tabs,
    TextField,
    Toolbar,
    Tooltip,
    Typography
} from '@mui/material';
import {
    AutorenewSharp as AutoUpdateIcon,
    Logout as LogoutIcon,
    Refresh as RefreshIcon,
    Settings as SettingsIcon,
    Speed as SpeedIcon
} from '@mui/icons-material';
import PropTypes from "prop-types";
import { getConnections } from './api/getConnections';
import { getRouterInfo } from './api/getRouterInfo';
import NetworkConfig from './NetworkConfig.jsx';
import IntranetConfig from './IntranetConfig.jsx';
import UserConfig from './UserConfig.jsx';
import BlackList from './BlackList.jsx';
import APIGuide from './APIGuide.jsx';
function Dashboard({ setIsLoggedIn }) {
    const [currentTab, setCurrentTab] = useState(0);

    const handleTabChange = (event, newValue) => {
        setCurrentTab(newValue);
    };

    const handleLogout = () => {
        setIsLoggedIn(false);
    };

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static">
            <Toolbar>
                <Typography 
                    variant="h6" 
                    component="div" 
                    sx={{ flexGrow: 1 }}
                >
                    라우터 관리 시스템
                </Typography>

                <Button
                    color="inherit"
                    onClick={() => setCurrentTab(5)}
                    sx={{
                        textDecoration: 'none',
                        color: 'inherit',
                        marginRight: 2,
                    }}
                >
                    API
                </Button>

                <Button
                    color="inherit"
                    onClick={handleLogout}
                    startIcon={<LogoutIcon />}
                >
                    로그아웃
                </Button>
            </Toolbar>
                <Tabs
                    value={currentTab}
                    onChange={handleTabChange}
                    sx={{ backgroundColor: 'primary.dark' }}
                >
                    <Tab icon={<SpeedIcon />} label="정보" />
                    <Tab icon={<SettingsIcon />} label="네트워크 기본 설정" />
                    <Tab icon={<SettingsIcon />} label="내부 네트워크 설정" />
                    <Tab icon={<SettingsIcon />} label="제한 목록" />
                    <Tab icon={<SettingsIcon />} label="관리자 설정" />
                    <Tab icon={<SettingsIcon />} label="API 가이드" />
                </Tabs>
            </AppBar>

            <Box sx={{ p: 4 }}>
                {currentTab === 0 && <InfoPanel />}
                {currentTab === 1 && <NetworkConfig />}
                {currentTab === 2 && <IntranetConfig />}
                {currentTab === 3 && <BlackList />}
                {currentTab === 4 && <UserConfig />}
                {currentTab === 5 && <APIGuide />}
            </Box>
        </Box>
    );
};

function InfoPanel() {
    const [routerInfo, setRouterInfo] = useState([]);
    const [connectionLog, setConnectionLog] = useState([]);
    const [isAutoUpdate, setIsAutoUpdate] = useState(true);
    const [updateError, setUpdateError] = useState(null);
    const [selectedIP, setSelectedIP] = useState(null);
    const [ipInfo, setIpInfo] = useState([]);
    const [logIndex, setLogIndex] = useState({});
    const fetchRouterInfo = async () => {
        try {
            const [routerData, connectionsData] = await Promise.all([
                getRouterInfo(),
                getConnections()
            ]);
            if (routerData) {
                setRouterInfo(routerData);
                setUpdateError(null);
            }

            if (Array.isArray(connectionsData)) setConnectionLog(connectionsData);
            else setConnectionLog([]);
        
        } catch (error) {
            setUpdateError('라우터 정보 업데이트 중 오류가 발생했습니다.');
            console.error('Router info update error:', error);
        }
    };

    useEffect(() => {
        fetchRouterInfo();
    }, []);

    useEffect(() => {
        if (!isAutoUpdate) return;
        
        const intervalId = setInterval(fetchRouterInfo, 5000);
        return () => clearInterval(intervalId);
    }, [isAutoUpdate]);

    useEffect(() => {
        if (selectedIP) {
            setIpInfo(connectionLog.filter(log => 
                log.source_ip === selectedIP || log.dest_ip === selectedIP
            ));
        }
    }, [selectedIP, connectionLog]);

    return (
        <Stack spacing={3}>
            <Box display="flex" alignItems="center" justifyContent="flex-end" mb={2}>
                <FormControlLabel
                    control={
                        <Switch
                            checked={isAutoUpdate}
                            onChange={() => setIsAutoUpdate(!isAutoUpdate)}
                            color="primary"
                        />
                    }
                    label="자동 업데이트"
                />
                <Tooltip title="수동 업데이트">
                    <IconButton onClick={fetchRouterInfo}>
                        <RefreshIcon />
                    </IconButton>
                </Tooltip>
            </Box>

            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        패킷 통신 로그
                    </Typography>
                    <TableContainer
                        component={Paper}
                        sx={{
                            maxHeight: '500', // 제한된 높이 설정
                            overflowY: 'auto', // 수직 스크롤 활성화
                        }}
                    >
                        <Table stickyHeader>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold' }}>프로토콜</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>출발 IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>도착 IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>출발 포트</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>도착 포트</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>패킷 수</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>바이트 크기</TableCell>
                                </TableRow>
                            </TableHead>
                            
                            <TableBody>
                                {updateError ? (
                                    <TableRow>
                                    <TableCell colSpan={7} align="center">
                                        <Typography color="error" variant="body1">
                                        {updateError}
                                        </Typography>
                                    </TableCell>
                                    </TableRow>
                                ) : connectionLog.map((log, index) => {
                                    
                                    // const key = `${log.protocol}${log.source_ip}${log.dest_ip}${log.source_port}`;
                                    // // prev -> next 형태를 명시적으로 반환
                                    // setLogIndex(prev => ({
                                    // ...prev,
                                    // [key]: prev[key] !== undefined ? `${prev[key]},${index}` : `${index}`,
                                    // }));

                                    return log.source_ip === log.dest_ip ? null : (
                                        <TableRow key={index}>
                                        <TableCell>{log.protocol}</TableCell>
                                        <TableCell
                                            sx={{ cursor: 'pointer', color: 'blue' }}
                                            onClick={() => setSelectedIP(log.source_ip)}
                                        >
                                            {log.source_ip}
                                        </TableCell>
                                        <TableCell
                                            sx={{ cursor: 'pointer', color: 'blue' }}
                                            onClick={() => setSelectedIP(log.dest_ip)}
                                        >
                                            {log.dest_ip}
                                        </TableCell>
                                        <TableCell>{log.source_port}</TableCell>
                                        <TableCell>{log.dest_port}</TableCell>
                                        <TableCell>{log.packet_count}</TableCell>
                                        <TableCell>{log.byte_count}</TableCell>
                                        </TableRow>
                                    );
                                })}
                            </TableBody>

                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>

            {selectedIP && (
                <Card>
                    <CardContent>
                        <Typography variant="h6">
                            선택된 IP: {selectedIP}의 연결 정보
                        </Typography>
                        <TableContainer 
                            component={Paper} 
                            sx={{ 
                                maxHeight: 300, 
                                overflowY: 'auto' 
                            }}
                        >
                            <Table stickyHeader>
                                <TableBody>
                                    {ipInfo.map((log, index) => (
                                        <TableRow key={index}>
                                            <TableCell>{log.source_ip === selectedIP ? log.dest_ip : log.source_ip}</TableCell>
                                            <TableCell>{log.source_ip === selectedIP ? log.source_port : log.dest_port}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </CardContent>
                </Card>
            )}            
        
            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        라우터 정보
                    </Typography>
                    <TableContainer component={Paper}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold' }}>항목</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>값</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {routerInfo.map((row) => (
                                    <TableRow key={row.name}>
                                        <TableCell>{row.name}</TableCell>
                                        <TableCell>{row.value}</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
        </Stack>
    );
}

Dashboard.propTypes = {
    setIsLoggedIn: PropTypes.func.isRequired
}

export default Dashboard;