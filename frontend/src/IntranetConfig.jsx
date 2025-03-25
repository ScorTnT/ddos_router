import { useEffect, useState } from 'react';
import {
    Button,
    Card,
    CardContent,
    TextField,
    Stack,
    TableContainer,
    Table,
    TableHead,
    TableBody,
    TableRow,
    TableCell,
    Paper,
    Typography,

} from "@mui/material";
import { LoadIntranetConfig, SaveIntranetConfig } from './api/intranetConfig';

function IntranetConfig() {
    const [ipAddress, setIpAddress] = useState("192.168.0.1");
    const [subnetMask, setSubnetMask] = useState("255.255.255.0");
    return (
        <>
        <Card>
            <CardContent>
                <Stack spacing={3}>
                    <Typography variant="h6" gutterBottom>
                        내부 네트워크 설정
                    </Typography>
                    
                    <TextField
                        label="내부 IP 주소"
                        variant="outlined"
                        value={ipAddress}
                        onChange={(e) => setIpAddress(e.target.value)}
                        fullWidth
                    />
                    
                    <TextField
                        label="내부 서브넷 마스크"
                        variant="outlined"
                        value={subnetMask}
                        onChange={(e) => setSubnetMask(e.target.value)}
                        fullWidth
                    />

                    <Typography variant="body2" color="textSecondary">
                        (내부 DHCP 설정은 추후 추가 예정)
                    </Typography>
                    
                    <Button
                        variant="contained"
                        color="primary"
                        sx={{ mt: 2 }}
                    >
                        설정 저장
                    </Button>
                </Stack>
            </CardContent>
        </Card>
        <IntranetIP></IntranetIP>
        </>
    );
}
function IntranetIP() {
    const [intrnetConnection, setIntranetConnection] = useState([]);

    const fetchIntranetConfig = async () => {
        const [intranetData] = await Promise.all([
            LoadIntranetConfig()
        ]);
        if(intranetData){
            setIntranetConnection(intranetData);
        }
    }

    useEffect(() => {
        fetchIntranetConfig();  
    },[]);

    return <>
    <Card>
        <CardContent>
            <Typography variant="h6" gutterBottom>
                내부 IP 목록
            </Typography>
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {intrnetConnection.map((row) => (
                            <TableRow key={row.name}>
                                <TableCell>{row.name}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </CardContent>
    </Card>
</>;
}

export default IntranetConfig;