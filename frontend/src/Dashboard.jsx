import {useEffect, useState} from 'react';
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
import {getConnections} from './api/getConnections';
import { getRouterInfo } from './api/getRouterInfo';
import NetworkConfig from './NetworkConfig.jsx';
import IntranetConfig from './IntranetConfig.jsx';
import UserConfig from './UserConfig.jsx';
import { getHardware } from './api/hardConfig.js';
function Dashboard({setIsLoggedIn}) {
    const [currentTab, setCurrentTab] = useState(0);

    const handleTabChange = (event, newValue) => {
        setCurrentTab(newValue);
    };

    const handleLogout = () => {
        setIsLoggedIn(false);
    };

    return (
        <Box sx={{flexGrow: 1}}>
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" component="div" sx={{flexGrow: 1}}>
                        라우터 관리 시스템
                    </Typography>
                    <Button
                        color="inherit"
                        onClick={handleLogout}
                        startIcon={<LogoutIcon/>}
                    >
                        로그아웃
                    </Button>
                </Toolbar>
                <Tabs
                    value={currentTab}
                    onChange={handleTabChange}
                    sx={{backgroundColor: 'primary.dark'}}
                >
                    <Tab icon={<SpeedIcon/>} label="정보"/>
                    <Tab icon={<SettingsIcon />} label="네트워크 기본 설정" />
                    <Tab icon={<SettingsIcon />} label="내부 네트워크 설정" />
                    <Tab icon={<SettingsIcon />} label="관리자 설정" />
                </Tabs>
            </AppBar>

            <Box sx={{ p: 3 }}>
                {currentTab === 0 && <InfoPanel />}
                {currentTab === 1 && <NetworkConfig />}
                {currentTab === 2 && <IntranetConfig />}
                {currentTab === 3 && <UserConfig />}
            </Box>
        </Box>
    );
};

function InfoPanel() {
    const [routerInfo, setRouterInfo] = useState([]);
    const [connectionLog, setConnectionLog] = useState('연결 정보를 불러오는 중...');
    const [isAutoUpdate, setIsAutoUpdate] = useState(true);
    const [updateError, setUpdateError] = useState(null);

    const fetchRouterInfo = async () => {
        try {
            const routerData = await getRouterInfo();
            if (routerData) {
                setRouterInfo([
                    { name: 'MAC 주소', value: routerData.mac_address || '알 수 없음' },
                    { name: '모델명', value: routerData.model_name || '알 수 없음' },
                    { name: '펌웨어 버전', value: routerData.firmware_version || '알 수 없음' },
                    { name: 'CPU 사용률', value: `${routerData.cpu_usage || 0}%` },
                    { name: '메모리 사용률', value: `${routerData.memory_usage || 0}%` },
                    { name: '가동 시간', value: routerData.uptime || '알 수 없음' },
                    { name: '연결된 기기 수', value: `${routerData.connected_devices || 0}대` },
                    { name: '현재 다운로드 속도', value: `${routerData.download_speed || 0}Mbps` },
                    { name: '현재 업로드 속도', value: `${routerData.upload_speed || 0}Mbps` },
                ]);
                setConnectionLog(routerData.packet_logs || '로그 정보 없음');
                setUpdateError(null);
            } else {
                setUpdateError('데이터를 불러올 수 없습니다.');
            }
        } catch (error) {
            setUpdateError('라우터 정보 업데이트 중 오류가 발생했습니다.');
            console.error('Router info update error:', error);
        }
    };
    
    useEffect(() => {
        fetchRouterInfo();
    }, []);

    useEffect(() => {
        let intervalId;
        if (isAutoUpdate) {
            intervalId = setInterval(fetchRouterInfo, 5000);
        }
        return () => {
            if (intervalId) {
                clearInterval(intervalId);
            }
        };
    }, [isAutoUpdate]);

    const routerInfo = [
        {name: 'MAC 주소', value: '00:11:22:33:44:55'},
        {name: '모델명', value: 'RT-AC86U'},
        {name: '펌웨어 버전', value: '3.0.0.4.386_45899'},
        {name: 'CPU 사용률', value: '25%'},
        {name: '메모리 사용률', value: '45%'},
        {name: '가동 시간', value: '10일 5시간 30분'},
        {name: '연결된 기기 수', value: '8대'},
        {name: '현재 다운로드 속도', value: '50Mbps'},
        {name: '현재 업로드 속도', value: '20Mbps'}
    ];

    return (
        <Stack spacing={3}>
            <Card>
                <CardContent>
                    <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
                        <Typography variant="h6">패킷 통신 로그</Typography>
                        <Box>
                            <Tooltip title={isAutoUpdate ? "자동 업데이트 중지" : "자동 업데이트 시작"}>
                                <IconButton
                                    onClick={() => setIsAutoUpdate(!isAutoUpdate)}
                                    color={isAutoUpdate ? "primary" : "default"}
                                >
                                    <AutoUpdateIcon />
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="수동 업데이트">
                                <IconButton onClick={fetchRouterInfo}>
                                    <RefreshIcon />
                                </IconButton>
                            </Tooltip>
                        </Box>
                    </Box>

                    <TextField
                        fullWidth
                        multiline
                        rows={10}
                        variant="outlined"
                        value={updateError || connectionLog}
                        InputProps={{
                            readOnly: true,
                            style: {
                                fontFamily: 'monospace',
                                fontSize: '0.875rem',
                            },
                        }}
                        error={!!updateError}
                        sx={{
                            backgroundColor: 'background.paper',
                            '& .MuiOutlinedInput-root': {
                                '& fieldset': {
                                    borderColor: updateError ? 'error.main' : 'rgba(0, 0, 0, 0.23)',
                                },
                            },
                        }}
                    />
                    {updateError && (
                        <Typography color="error" variant="caption" sx={{ mt: 1 }}>
                            {updateError}
                        </Typography>
                    )}
                </CardContent>
            </Card>

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