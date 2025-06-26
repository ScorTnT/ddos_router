import { useEffect, useState } from "react";
import {
    Button,
    Box,
    Card,
    CardContent,
    TextField,
    Stack,
    TableContainer,
    Table,
    TableHead,
    TableBody,
    TableCell,
    Paper,
    Typography,
    Tooltip,
    IconButton,
    TableRow,
    Checkbox
} from "@mui/material";
import { getBlockedIP, blockIP, unblockIP } from './api/getBlockedIP';
function BlackList(){
    const [protectionLog, setProtectionLog] = useState([]);
    // 리스트 관리
    const [blockedIP, setBlockedIP] = useState([]);
    const [blackList, setBlackList] = useState([]);
    const [whiteList, setWhiteList] = useState([]);
    // 선택된 리스트 관리
    const [selectedBlockedIP, setSelectedBlockedIP] = useState([]);
    const [selectedBlackList, setSelectedBlackList] = useState([]);
    const [selectedWhiteList, setSelectedWhiteList] = useState([]);

    const [isAutoUpdate, setIsAutoUpdate] = useState(true);
    const [updateError, setUpdateError] = useState(null);

    const FetchProtectionLog = async () => {
        try {
            const protectionData = await getBlockedIP();
            if (Array.isArray(protectionData)) setProtectionLog(protectionData);
            else setProtectionLog([]);

        } catch (error) {
            setUpdateError('제한 목록 업데이트 중 오류가 발생했습니다.');
            console.error('BlackList update error:', error);
        }
    };
    // 토글 스위찌 체크박스
    const CheckBlockedIP = (ip) => {
        setSelectedBlockedIP((prev) =>
            prev.includes(ip) ? prev.filter((item => item !== ip)) : [...prev, ip]
        );
    };
    const CheckBlackList = (ip) => {
        setSelectedBlackList((prev) =>
            prev.includes(ip) ? prev.filter((item => item !== ip)) : [...prev, ip]
        );
    };
    const CheckWhiteList = (ip) => {
        setSelectedWhiteList((prev) =>
            prev.includes(ip) ? prev.filter((item => item !== ip)) : [...prev, ip]
        );
    };
    // 버튼
    const ButtonBlock_White = () => {
        const whiteIP = blockedIP.filter(item => selectedBlockedIP.includes(item.ip));
        setWhiteList((prev) => [...prev, ...whiteIP])
        setBlockedIP((prev) => prev.filter(item => !selectedBlockedIP.includes(item.ip)));
        setSelectedBlockedIP([]);
    };
    const ButtonBlock_Black = () => {
        const blackIP = blockedIP.filter(item => selectedBlockedIP.includes(item.ip));
        setBlackList((prev) => [...prev, ...blackIP])
        setBlockedIP((prev) => prev.filter(item => !selectedBlockedIP.includes(item.ip)));
        setSelectedBlockedIP([]);
    }
    const ButtonBlack_White = () => {
        const whiteIP = blackList.filter(item => selectedBlackList.includes(item.ip));
        setWhiteList((prev) => [...prev, ...whiteIP])
        setBlockedIP((prev) => prev.filter(item => !selectedBlackList.includes(item.ip)));
        setSelectedBlackList([]);
    }
    const ButtonWhite_Black = () => {
        const blackIP = whiteList.filter(item => selectedWhiteList.includes(item.ip));
        setWhiteList((prev) => [...prev, ...whiteIP])
        setBlockedIP((prev) => prev.filter(item => !selectedWhiteList.includes(item.ip)));
        setSelectedWhiteList([]);
    }
    // fetch
    useEffect(() => {
        FetchProtectionLog();
    }, []);
    // update
    useEffect(() => {
        if (!isAutoUpdate) return;
        const intervalId = setInterval(fetchProtectionLog, 5000);
        return () => clearInterval(intervalId);       
    }, [isAutoUpdate]);
    // lists
    useEffect(() => {
        // api 답변 오류
        if (!Array.isArray(protectionLog)) return;
        // blackList 빈 경우
        if (blackList.length === 0) {
            setBlockedIP(protectionLog);
            return;
        }
        // 정상적인 경우 (새로 생긴 ip만만 BlockedIP 리스트에 추가)
        const pLog = [...protectionLog];
        const pastLogIp = [
            ...blockedIP.map(item => item.ip),
            ...blackList.map(item => item.ip),
        ];
        const newLog = pLog.filter(item => !pastLogIp.includes(item.ip));       
        setBlockedIP(prev => [...prev, ...newLog]);
    }, [protectionLog]);
    
    return <>
        <Stack spacing={3}>
            <Card>
                <CardContent>
                    <Box sx={{
                        display : "flex",
                        justifyContent : "space-between",
                        alignItems : "center",
                        mb : 2
                    }}>
                        <Typography variant="h6" gutterBottom>
                            차단 IP 목록 (Dropped IP)
                        </Typography>

                        <Stack direction="row" spacing={2}>
                            <Button 
                                variant="contained" 
                                color="success" 
                                onClick={ButtonBlock_White}
                                disabled={selectedBlockedIP.length===0}
                            >
                                허용
                            </Button>
                            <Button 
                                variant="contained" 
                                color="error" 
                                onClick={ButtonBlock_Black}
                                disabled={selectedBlockedIP.length===0}
                            >
                                차단
                            </Button>
                        </Stack>
                    </Box>

                    <TableContainer 
                        component={Paper}
                        sx={{
                            overflowY : 'auto',
                        }}
                    >
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>차단된 시각</TableCell>
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
                                ) : blockedIP.map((pLog, index) => {
                                    <TableRow key={index} selected={selectedBlockedIP.includes(pLog.ip)}>
                                        <TableCell padding="checkbox">
                                            <Checkbox
                                                checked={selectedBlockedIP.includes(pLog.ip)}
                                                onChange={() => CheckBlockedIP(pLog.ip)}
                                            />
                                        </TableCell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        <TableCell>{formatTime(pLog.expire_at)}</TableCell>
                                    </TableRow>
                                })}
                            </TableBody>

                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>

            <Card>
                <CardContent>
                    <Box sx={{
                        display : "flex",
                        justifyContent : "space-between",
                        alignItems : "center",
                        mb : 2
                    }}>
                        <Typography variant="h6" gutterBottom>
                            영구 차단 IP 목록 (IP BlackList)
                        </Typography>
                        <Stack direction="row" spacing={2}>
                            <Button 
                                variant="contained" 
                                color="success" 
                                onClick={ButtonBlack_White}
                                disabled={selectedBlackList.length===0}
                            >
                                허용
                            </Button>
                        </Stack>
                    </Box>

                    <TableContainer 
                        component={Paper}
                        sx={{
                            overflowY : 'auto',
                        }}
                    >
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>차단된 시각</TableCell>
                                </TableRow>
                            </TableHead>

                            <TableBody>
                                {blackList.length = 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={7} align="center">
                                            <Typography color="error" variant="body1">
                                                차단된 IP가 존재하지 않습니다.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>                                    
                                ) : blackList.map((pLog, index) => {
                                    <TableRow key={index} selected={selectedBlackList.includes(pLog.ip)}>
                                        <Tablecell padding="checkbox">
                                            <Checkbox
                                                checked={selectedBlackList.includes(pLog.ip)}
                                                onChange={() => CheckBlackList(pLog.ip)}
                                            />
                                        </Tablecell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        <TableCell>{formatTime(pLog.expire_at)}</TableCell>
                                    </TableRow>
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>                    
                </CardContent>
            </Card>

            <Card>
                <CardContent>
                    <Box sx={{
                        display : "flex",
                        justifyContent : "space-between",
                        alignItems : "center",
                        mb : 2
                    }}>
                        <Typography variant="h6" gutterBottom>
                            허용 IP 목록 (IP WhiteList)
                        </Typography>
                        <Stack direction="row" spacing={2}>
                            <Button 
                                variant="contained" 
                                color="success" 
                                onClick={ButtonWhite_Black}
                                disabled={selectedWhiteList.length===0}
                            >
                                차단
                            </Button>
                        </Stack>
                    </Box>

                    <TableContainer 
                        component={Paper}
                        sx={{
                            overflowY : 'auto',
                        }}
                    >
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                                    {/*<TableCell sx={{ fontWeight: 'bold' }}>허용된 시각</TableCell>*/}
                                </TableRow>
                            </TableHead>

                            <TableBody>
                                {whiteList.length = 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={7} align="center">
                                            <Typography color="error" variant="body1">
                                                허용된 IP가 존재하지 않습니다.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>                                    
                                ) : whiteList.map((pLog, index) => {
                                    <TableRow key={index}>
                                        <TableCell padding="checkbox">
                                            <Checkbox
                                                checked={selectedWhiteList.includes(pLog.ip)}
                                                onChange={() => CheckWhiteList(pLog.ip)}
                                            />
                                        </TableCell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        {/*<TableCell>{formatTime(pLog.expire_at)}</TableCell>*/}
                                    </TableRow>
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>                    
                </CardContent>
            </Card>
        </Stack>
    </>
}

export default BlackList;