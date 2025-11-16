import { useEffect, useState } from "react";
import {
    Button,
    Box,
    Card,
    CardContent,
    Stack,
    TableContainer,
    Table,
    TableHead,
    TableBody,
    TableCell,
    Paper,
    Typography,
    TableRow,
    Checkbox
} from "@mui/material";
import api from './api.js';

const IP_STATUS = {
    BLACKLIST: 'blacklist',
    WHITELIST: 'whitelist',
};

// 시간 포맷 함수는 이제 사용되지 않으므로 제거하거나 그대로 둡니다 (여기서는 제거합니다).
/*
const formatTime = (timestamp) => {
    if (!timestamp) return '-';
    const date = new Date(timestamp * 1000); 
    return date.toLocaleString();
};
*/

function BlackList() {
    const [ipList, setIpList] = useState([]);
    const [selectedIps, setSelectedIps] = useState(new Set());
    const [updateError, setUpdateError] = useState(null);

    const fetchProtectionLog = async () => {
        try {
            const rawBlackList = await api.getProtection();
            const rawWhiteList = await api.getWhiteList();
            
            const blackListItems = rawBlackList.map(ipObj => ({
                ...ipObj,
                status : IP_STATUS.BLACKLIST,
                isSelected : selectedIps.has(ipObj.ip),
            }));

            const whiteListItems = rawWhiteList.map(ipObj => {
                const ip = typeof ipObj === 'string' ? ipObj : ipObj.ip;
                return {
                    ip: ip,
                    // 시간 필드 제거
                    status : IP_STATUS.WHITELIST,
                    isSelected : selectedIps.has(ip),
                };
            });
            
            setIpList([...blackListItems, ...whiteListItems]);
            setUpdateError(null);
        } catch (error) {
            setUpdateError('제한 목록을 불러오는 중 오류가 발생했습니다.');
            console.error('Protection logs fetch error:', error);
        }
    };

    const handleCheck = (ip) => {
        setSelectedIps(prev => {
            const newSelected = new Set(prev);
            if (newSelected.has(ip)) {
                newSelected.delete(ip);
            } else {
                newSelected.add(ip);
            }
            return newSelected;
        });

        setIpList((prev) =>
            prev.map((item) =>
                item.ip === ip ? { ...item, isSelected: !item.isSelected } : item
            )
        );
    };

    const handleWhite = async () => {
        const selectedBlackIps = ipList
            .filter(item => selectedIps.has(item.ip) && item.status === IP_STATUS.BLACKLIST)
            .map(item => item.ip);
            
        if (selectedBlackIps.length === 0) {
            alert('허용할 IP가 선택되지 않았습니다.');
            return;
        }

        try {
            for (const ip of selectedBlackIps) {
                const addWhiteListRes = await api.addWhiteList(ip);
                console.log(`[API LOG] ${ip} - addWhiteList 응답:`, addWhiteListRes);

                const unblockIPRes = await api.unblockIP(ip);
                console.log(`[API LOG] ${ip} - unblockIP 응답:`, unblockIPRes);
            }
            
            await fetchProtectionLog();
            setSelectedIps(new Set());

        } catch (error) {
            console.error('IP 허용 오류', error);
            alert('IP 허용 중 오류가 발생했습니다.');
        }
    };

    const handleBlack = async () => {
        const selectedWhiteIps = ipList
            .filter(item => selectedIps.has(item.ip) && item.status === IP_STATUS.WHITELIST)
            .map(item => item.ip);

        if (selectedWhiteIps.length === 0) {
            alert('차단할 IP가 선택되지 않았습니다.');
            return;
        }

        try {
            for (const ip of selectedWhiteIps) {
                // 1. IP 차단 API 호출 및 응답 로깅
                const blockIPRes = await api.blockIP(ip);
                console.log(`[API LOG] ${ip} - blockIP 응답:`, blockIPRes);

                // 2. WhiteList 제거 API 호출 및 응답 로깅
                const removeWhiteListRes = await api.removeWhiteList(ip);
                console.log(`[API LOG] ${ip} - removeWhiteList 응답:`, removeWhiteListRes);
            }
            
            // API 호출 후 서버에서 최신 상태를 다시 가져와 일괄 업데이트
            await fetchProtectionLog();
            setSelectedIps(new Set());

        } catch (error) {
            // ★ 중요한 디버깅 지점: 에러가 발생하면 어떤 IP의 어떤 API에서 문제가 발생했는지 확인
            console.error('IP 차단 오류 - 상세 에러 객체:', error);
            alert(`IP 차단 중 오류가 발생했습니다: ${error.message}`);
        }
    };

    useEffect(() => {
        fetchProtectionLog();
        const intervalId = setInterval(fetchProtectionLog, 5000);
        return () => clearInterval(intervalId);
    }, []);
    
    const blackList = ipList.filter(item => item.status === IP_STATUS.BLACKLIST);
    const whiteList = ipList.filter(item => item.status === IP_STATUS.WHITELIST);

    const selectedBlackList = blackList.filter(item => selectedIps.has(item.ip)).length;
    const selectedWhiteList = whiteList.filter(item => selectedIps.has(item.ip)).length;

    return (
        <Stack spacing={3}>
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
                                onClick={handleWhite}
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
                                    {/* <TableCell sx={{ fontWeight: 'bold' }}>차단 해제 시각</TableCell> */}
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
                                ) : blackList.map((pLog) => (
                                    <TableRow key={pLog.ip} selected={selectedIps.has(pLog.ip)}>
                                        <TableCell padding="checkbox">
                                            <Checkbox
                                                checked={selectedIps.has(pLog.ip)}
                                                onChange={() => handleCheck(pLog.ip)}
                                            />
                                        </TableCell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        {/* <TableCell>{formatTime(pLog.expire_at)}</TableCell> */}
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer> 
                </CardContent>
            </Card>

            <Card>
                <CardContent>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
                        <Typography variant="h6" gutterBottom>
                            허용 IP 목록 (IP WhiteList)
                        </Typography>
                        <Stack direction="row" spacing={2}>
                            <Button 
                                variant="contained" 
                                color="error" 
                                onClick={handleBlack}
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
                                    {/* <TableCell sx={{ fontWeight: 'bold' }}>추가된 시각</TableCell> */}
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
                                ) : whiteList.map((pLog) => (
                                    <TableRow key={pLog.ip} selected={selectedIps.has(pLog.ip)}>
                                        <TableCell padding="checkbox">
                                            <Checkbox
                                                checked={selectedIps.has(pLog.ip)}
                                                onChange={() => handleCheck(pLog.ip)}
                                            />
                                        </TableCell>
                                        <TableCell>{pLog.ip}</TableCell>
                                        {/* <TableCell>{formatTime(pLog.added_at) || '-'}</TableCell> */}
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