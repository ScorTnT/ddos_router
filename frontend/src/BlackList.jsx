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

    /*
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
        const tempBlocked = [];
        const tempBlack = [];
        
        protectionLog.forEach(pLog => {
            const pLogIp = pLog.ip;
            if (blackList.some(item => item.ip === pLogIp)) {
                tempBlack.push();
            } else if (whiteList.includes(ip)) {
                tempWhite.push(ipObj);
            } else {
                tempBlocked.push(ipObj);
            }
        });

        setBlockedIP(tempBlocked);
        // setBlackList(tempBlack.map(i => i.ip)); // <- 필요하면
        // setWhiteList(tempWhite.map(i => i.ip)); // <- 필요하면
    }, [apiList, blackList, whiteList]);
    
    */

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
                                <TableRow>
                                    <TableCell> api 답변에 따라 수정 예정 </TableCell>
                                    <TableCell> api 답변에 따라 수정 예정 </TableCell>
                                </TableRow>
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
                                    <TableCell sx={{ fontWeight: 'bold' }}>추가된 시각</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                <TableRow>
                                    <TableCell>     </TableCell>
                                    <TableCell>     </TableCell>
                                </TableRow>
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
                                <TableRow>
                                    <TableCell>     </TableCell>
                                    <TableCell>     </TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>
                    </TableContainer>                    
                </CardContent>
            </Card>
        </Stack>
    </>
}

export default BlackList;