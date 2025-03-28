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
    Typography
} from "@mui/material";
import { LoadIntranetConfig, SaveIntranetConfig } from './api/intranetConfig';
import { getArpInfo } from './api/arpConfig';
function IntranetConfig() {
    const [ipAddress, setIpAddress] = useState("192.168.0.1");
    const [subnetMask, setSubnetMask] = useState("255.255.255.0");
    function fetchIntranetConfig() {
        LoadIntranetConfig().then((data) => {
            if (data) {
                setIpAddress(data.IPAddress || "");
                setSubnetMask(data.Netmask || "");
            }
        }).catch((error) => {
            console.error("Error loading intranet config:", error);
        });
    }
    useEffect(() => {
        fetchIntranetConfig();
    }
    , []);

    return (
        <Stack spacing={2}>
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
        </Stack>
    );
}
function IntranetIP() {
    const [intranetConnection, setIntranetConnection] = useState([]);

    const fetchIntranetConfig = async () => {
        try {
            const intranetData = await getArpInfo();

            if(intranetData){
                setIntranetConnection(intranetData);
            }
        } catch (error) {
            console.error('Error fetching intranet data:', error);
        }
    }

    useEffect(() => {
        fetchIntranetConfig();  
    },[]);

    return (
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
                            {intranetConnection.map((row) => (
                                <TableRow key={row.ip}>
                                    <TableCell>{row.ip}</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </CardContent>
        </Card>
    );
}

export default IntranetConfig;