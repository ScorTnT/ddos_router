import React, {useEffect, useState} from 'react';
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
import {getConnections} from './api/getConnections';

const Dashboard = ({setIsLoggedIn}) => {
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
                    <Tab icon={<SettingsIcon/>} label="설정"/>
                </Tabs>
            </AppBar>

            <Box sx={{p: 3}}>
                {currentTab === 0 ? (
                    <InfoPanel/>
                ) : (
                    <SettingsPanel/>
                )}
            </Box>
        </Box>
    );
};

const InfoPanel = () => {
    const [connectionLog, setConnectionLog] = useState('연결 정보를 불러오는 중...');
    const [isAutoUpdate, setIsAutoUpdate] = useState(true);
    const [updateError, setUpdateError] = useState(null);

    const fetchConnectionData = async () => {
        try {
            const data = await getConnections();
            if (data) {
                setConnectionLog(data);
                setUpdateError(null);
            } else {
                setUpdateError('데이터를 불러올 수 없습니다.');
            }
        } catch (error) {
            setUpdateError('연결 정보 업데이트 중 오류가 발생했습니다.');
            console.error('Connection update error:', error);
        }
    };

    useEffect(() => {
        fetchConnectionData();
    }, []);

    useEffect(() => {
        let intervalId;
        if (isAutoUpdate) {
            intervalId = setInterval(fetchConnectionData, 5000);
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
                        <Typography variant="h6">
                            패킷 통신 로그
                        </Typography>
                        <Box>
                            <Tooltip title={isAutoUpdate ? "자동 업데이트 중지" : "자동 업데이트 시작"}>
                                <IconButton
                                    onClick={() => setIsAutoUpdate(!isAutoUpdate)}
                                    color={isAutoUpdate ? "primary" : "default"}
                                >
                                    <AutoUpdateIcon/>
                                </IconButton>
                            </Tooltip>
                            <Tooltip title="수동 업데이트">
                                <IconButton onClick={fetchConnectionData}>
                                    <RefreshIcon/>
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
                        inputProps={{
                            readOnly: true,
                            style: {
                                fontFamily: 'monospace',
                                fontSize: '0.875rem'
                            }
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
                        <Typography color="error" variant="caption" sx={{mt: 1}}>
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
                                    <TableCell sx={{fontWeight: 'bold'}}>항목</TableCell>
                                    <TableCell sx={{fontWeight: 'bold'}}>값</TableCell>
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
};

const SettingsPanel = () => {
    return (
        <Card>
            <CardContent>
                <Stack spacing={3}>
                    <Typography variant="h6" gutterBottom>
                        기본 설정
                    </Typography>

                    <FormControl fullWidth>
                        <InputLabel>무선 채널</InputLabel>
                        <Select
                            value={6}
                            label="무선 채널"
                        >
                            <MenuItem value={1}>채널 1 (2.4GHz)</MenuItem>
                            <MenuItem value={6}>채널 6 (2.4GHz)</MenuItem>
                            <MenuItem value={11}>채널 11 (2.4GHz)</MenuItem>
                        </Select>
                    </FormControl>

                    <FormControl fullWidth>
                        <InputLabel>보안 모드</InputLabel>
                        <Select
                            value="wpa3"
                            label="보안 모드"
                        >
                            <MenuItem value="wpa2">WPA2-PSK</MenuItem>
                            <MenuItem value="wpa3">WPA3-PSK</MenuItem>
                            <MenuItem value="mixed">WPA2/WPA3 혼합</MenuItem>
                        </Select>
                    </FormControl>

                    <FormControlLabel
                        control={<Switch defaultChecked/>}
                        label="게스트 네트워크 활성화"
                    />

                    <FormControlLabel
                        control={<Switch defaultChecked/>}
                        label="QoS 활성화"
                    />

                    <FormControlLabel
                        control={<Switch defaultChecked/>}
                        label="방화벽 활성화"
                    />

                    <Button
                        variant="contained"
                        color="primary"
                        sx={{mt: 2}}
                    >
                        설정 저장
                    </Button>
                </Stack>
            </CardContent>
        </Card>
    );
};

export default Dashboard;