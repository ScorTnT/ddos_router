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
import { Refresh as RefreshIcon } from '@mui/icons-material';
import api from './api.js';

const IP_STATUS = {
    BLACKLIST: 'blacklist',
    WHITELIST: 'whitelist',
};

const formatTime = (timestamp) => {
    if (!timestamp) return '-';
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
};

function BlackList() {
    const [ipList, setIpList] = useState([]);
    const [isAutoUpdate, setIsAutoUpdate] = useState(true);
    const [updateError, setUpdateError] = useState(null);

    const fetchProtectionLog = async () => {
        try {
            const data = await api.getProtection();          
            const rawIpList = data.map(ipObj => ({
                ...ipObj,
                status : IP_STATUS.BLACKLIST,
                isSelected : false,
            }));
            
            setIpList(rawIpList);
            setUpdateError(null);
        } catch (error) {
            setUpdateError('제한 목록을 불러오는 중 오류가 발생했습니다.');
            console.error('BlackList fetch error:', error);
        }
    };

    const handleCheck = (ip) => {
        setIpList((prev) =>
            prev.map((item) =>
                item.ip === ip ? { ...item, isSelected: !item.isSelected } : item
            )
        );
    };

    // Black -> White
    const handleWhite = async () => {
        const selectedIp = ipList.filter(item => item.isSelected && item.status === IP_STATUS.BLACKLIST);
        if (selectedIp.length === 0) {
            alert('선택된 ip가 존재하지 않습니다.');
            return;
        }
        try {
            for (const ipObj of selectedIp) {
                await api.unblockIP(ipObj.ip);
            }
            setIpList(prev =>
                prev.map(item =>
                    item.isSelected && item.status === IP_STATUS.BLACKLIST
                        ? { ...item, status: IP_STATUS.WHITELIST, isSelected: false }
                        : item
                )
            );
        } catch (error) {
            console.error('IP 허용 오류', error);
            alert('IP 허용 중 오류가 발생했습니다.');
        }
    };

    // White -> Black
    const handleBlack = async () => {
        const selectedIp = ipList.filter(item => item.isSelected && item.status === IP_STATUS.WHITELIST);
        if (selectedIp.length === 0) {
            alert('선택된 ip가 존재하지 않습니다.');
            return;
        }
        try {
            for (const ipObj of selectedIp) {
                await api.blockIP(ipObj.ip);
            }
            setIpList(prev =>
                prev.map(item =>
                    item.isSelected && item.status === IP_STATUS.WHITELIST
                        ? { ...item, status: IP_STATUS.BLACKLIST, isSelected: false }
                        : item
                )
            );
        } catch (error) {
            console.error('IP 차단 오류', error);
            alert('IP 차단 중 오류가 발생했습니다.');
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
    
    //const blockedIp = ipList.filter(item => item.status === IP_STATUS.BLOCKED);
    const blackList = ipList.filter(item => item.status === IP_STATUS.BLACKLIST);
    const whiteList = ipList.filter(item => item.status === IP_STATUS.WHITELIST);
    //const selectedBlockedIp = blockedIp.filter(item => item.isSelected).length;
    const selectedBlackList = blackList.filter(item => item.isSelected).length;
    const selectedWhiteList= whiteList.filter(item => item.isSelected).length;

    return (
        <Stack spacing={3}>
            {/* 차단 IP 목록 (IP BlackList) */}
            <Card>
                <CardContent>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
                        <Typography variant="h6" gutterBottom>
                            차단 IP 목록 (IP BlackList)
                        </Typography>
                        <Stack direction="row" spacing={2}>
                            <Button 
                                variant="contained" 
                                color="success" 
                                onClick={() => handleWhite(IP_STATUS.BLACKLIST)}
                                disabled={selectedBlackList === 0}
                            >
                                허용
                            </Button>
                        </Stack>
                    </Box>
                    <TableContainer component={Paper} sx={{ overflowY: 'auto' }}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell padding="checkbox" />
                                    <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>차단된 시각</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {updateError ? (
                                    <TableRow>
                                        <TableCell colSpan={2} align="center">
                                            <Typography color="error" variant="body1">
                                                {updateError}
                                            </Typography>
                                        </TableCell>
                                    </TableRow>                                    
                                ) : blackList.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={2} align="center">
                                            <Typography color="textSecondary" variant="body1">
                                                차단된 IP가 존재하지 않습니다.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                ) : blackList.map((pLog, index) => (
                                    <TableRow key={index} selected={pLog.isSelected}>
                                        <TableCell padding="checkbox">
                                            <Checkbox
                                                checked={pLog.isSelected}
                                                onChange={() => handleCheck(pLog.ip)}
                                            />
                                        </TableCell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        <TableCell>{formatTime(pLog.expire_at)}</TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer> 
                </CardContent>
            </Card>

            {/* 허용 IP 목록 (IP WhiteList) */}
            <Card>
                <CardContent>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
                        <Typography variant="h6" gutterBottom>
                            허용 IP 목록 (IP WhiteList)
                        </Typography>
                        <Stack direction="row" spacing={2}>
                            <Button 
                                variant="contained" 
                                color="success" 
                                onClick={() => handleBlack(IP_STATUS.WHITELIST)}
                                disabled={selectedWhiteList === 0}
                            >
                                차단
                            </Button>
                        </Stack>
                    </Box>
                    <TableContainer component={Paper} sx={{ overflowY: 'auto' }}>
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell padding="checkbox" />
                                    <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>차단된 시각</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {whiteList.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={2} align="center">
                                            <Typography color="textSecondary" variant="body1">
                                                허용된 IP가 존재하지 않습니다.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                ) : whiteList.map((pLog, index) => (
                                    <TableRow key={index} selected={pLog.isSelected}>
                                        <TableCell padding="checkbox">
                                            <Checkbox
                                                checked={pLog.isSelected}
                                                onChange={() => handleCheck(pLog.ip)}
                                            />
                                        </TableCell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        <TableCell>{formatTime(pLog.expire_at)}</TableCell>
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

export default BlackList;