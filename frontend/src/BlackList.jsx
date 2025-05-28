import { useEffect, useState } from "react";
import {
    Button,
    Box,
    Button,
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
    TableRow
} from "@mui/material";
import { getBlockedIP, blockIP, unblockIP } from './api/getBlockedIP';
function BlackList(){
    const [protectionLog, setProtectionLog] = useState([]);
    const [blockedIP, setBlockedIP] = useState([]);
    const [blackList, setBlackList] = useState([]);
    const [whiteList, setWhiteList] = useState([]);
    const [isAutoUpdate, setIsAutoUpdate] = useState(true);
    const [updateError, setUpdateError] = useState(null);

    const fetchProtectionLog = async () => {
        try {
            const protectionData = await getBlockedIP();
            if (Array.isArray(protectionData)) setProtectionLog(protectionData);
            else setProtectionLog([]);

        } catch (error) {
            setUpdateError('제한 목록 업데이트 중 오류가 발생했습니다.');
            console.error('BlackList update error:', error);
        }
    };

    useEffect(() => {
        fetchProtectionLog();
    }, []);

    useEffect(() => {
        if (!isAutoUpdate) return;
        const intervalId = setInterval(fetchProtectionLog, 5000);
        return () => clearInterval(intervalId);       
    }, [isAutoUpdate]);

    useEffect(() => {
        // if (!Array.isArray(protectionLog)) return;
        const tempBlack = [...protectionLog];
        const tempBlocked = [];
        
        for (let i = tempBlack.length - 1; i >= 0; i--) {
            const pLog = tempBlack[i];
            const pLogIp = pLog.ip
            if (!blackList.some(item => item.ip === pLogIp)) {
                tempBlocked.push(pLog);
                tempBlack.splice(i, 1);
            }
        }

        setBlockedIP(tempBlocked.reverse());
        setBlackList(tempBlack);
    }, [protectionLog]);

    useEffect(() => {

    })
    
    return <>
        <Stack spacing={3}>
            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        차단 IP 목록 (Dropped IP)
                    </Typography>
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
                                    <TableRow key={index}>
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
                    <Typography variant="h6" gutterBottom>
                        영구 차단 IP 목록 (IP BlackList)
                    </Typography>
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
                                    <TableRow key={index}>
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
                    <Typography variant="h6" gutterBottom>
                        허용 IP 목록 (IP WhiteList)
                    </Typography>
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
                                    <TableCell sx={{ fontWeight: 'bold' }}>허용된 시각</TableCell>
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
                                        <TableCell>{pLog.ip}</TableCell>
                                        <TableCell>{formatTime(pLog.expire_at)}</TableCell>
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