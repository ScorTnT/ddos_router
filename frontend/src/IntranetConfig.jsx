/*
import {useState} from 'react'
import {
    Button,
    Card,
    CardContent,
    FormControl,
    FormControlLabel,
    InputLabel,
    MenuItem,
    Select,
    Stack, Switch,
    Typography
} from "@mui/material";

function IntranetConfig()  {
    const IntranetStatus = "intranet";

    return (
        <Card>
            <CardContent>
                <Stack spacing={3}>
                    <Typography variant="h6" gutterBottom>
                        내부 네트워크
                    </Typography>
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
}

export default IntranetConfig;
*/

import { useState } from 'react';
import {
    Button,
    Card,
    CardContent,
    TextField,
    Stack,
    Typography,
} from "@mui/material";

function IntranetConfig() {
    const [ipAddress, setIpAddress] = useState("192.168.0.1");
    const [subnetMask, setSubnetMask] = useState("255.255.255.0");
    const [dhcp_, setDhcp] = useState("");
    return (
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
    );
}

export default IntranetConfig;
