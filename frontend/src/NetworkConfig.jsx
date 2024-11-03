import { useState } from 'react';
import {
    Box,
    Button,
    Card,
    CardContent,
    FormControlLabel,
    Radio,
    RadioGroup,
    TextField,
    Stack,
    Switch,
    Typography
} from "@mui/material";

function NetworkConfig() {
    const [connectionType, setConnectionType] = useState("dynamic"); // 동적 IP 방식 기본값
    const [ipAddress, setIpAddress] = useState("220.66.87.40");
    const [subnetMask, setSubnetMask] = useState("255.255.255.0");
    const [gateway, setGateway] = useState("220.66.87.2");
    const [primaryDNS, setPrimaryDNS] = useState("8.8.8.8");
    const [secondaryDNS, setSecondaryDNS] = useState("8.8.4.4");
    const [wanMacAddress, setWanMacAddress] = useState("ff:ff:ff:ff:ff:ff");
    const [mtu, setMtu] = useState("1500");
    const [manualDns, setManualDns] = useState(true);
    const [macAddressChange, setMacAddressChange] = useState(true);
    const [manualMtu, setManualMtu] = useState(false);

    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom>
                    인터넷 설정 정보
                </Typography>
                <Stack spacing={2}>
                    <Box display="flex" alignItems="center">
                        <Typography>IP 방식</Typography>
                        <RadioGroup
                            row
                            value={connectionType}
                            onChange={(e) => setConnectionType(e.target.value)}
                            sx={{ marginLeft: 2 }}
                        >
                            <FormControlLabel
                                value="dynamic"
                                control={<Radio />}
                                label="동적 IP 방식"
                            />
                            <FormControlLabel
                                value="static"
                                control={<Radio />}
                                label="고정 IP 방식"
                            />
                        </RadioGroup>
                    </Box>

                    <TextField
                        label="외부 IP 주소"
                        value={ipAddress}
                        onChange={(e) => setIpAddress(e.target.value)}
                        fullWidth
                        disabled={connectionType === "dynamic"} // 동적 IP 방식일 때 비활성화
                    />
                    <TextField
                        label="서브넷 마스크"
                        value={subnetMask}
                        onChange={(e) => setSubnetMask(e.target.value)}
                        fullWidth
                        disabled={connectionType === "dynamic"} // 동적 IP 방식일 때 비활성화
                    />
                    <TextField
                        label="기본 게이트웨이"
                        value={gateway}
                        onChange={(e) => setGateway(e.target.value)}
                        fullWidth
                        disabled={connectionType === "dynamic"} // 동적 IP 방식일 때 비활성화
                    />

                    <FormControlLabel
                        control={<Switch checked={manualDns} onChange={(e) => setManualDns(e.target.checked)} />}
                        label="DNS 주소 수동 입력"
                    />
                    <TextField
                        label="기본 DNS 주소"
                        value={primaryDNS}
                        onChange={(e) => setPrimaryDNS(e.target.value)}
                        fullWidth
                        disabled={!manualDns}  // 수동 DNS 해제 시 비활성화
                    />
                    <TextField
                        label="보조 DNS 주소"
                        value={secondaryDNS}
                        onChange={(e) => setSecondaryDNS(e.target.value)}
                        fullWidth
                        disabled={!manualDns}  // 수동 DNS 해제 시 비활성화
                    />

                    <FormControlLabel
                        control={<Switch checked={macAddressChange} onChange={(e) => setMacAddressChange(e.target.checked)} />}
                        label="MAC 주소 변경"
                    />
                    <TextField
                        label="WAN MAC 주소"
                        value={wanMacAddress}
                        onChange={(e) => setWanMacAddress(e.target.value)}
                        fullWidth
                        disabled={!macAddressChange}  // MAC 주소 변경 해제 시 비활성화
                    />

                    <FormControlLabel
                        control={<Switch checked={manualMtu} onChange={(e) => setManualMtu(e.target.checked)} />}
                        label="MTU 수동 입력"
                    />
                    <TextField
                        label="MTU"
                        value={mtu}
                        onChange={(e) => setMtu(e.target.value)}
                        fullWidth
                        disabled={!manualMtu}  // MTU 수동 입력 해제 시 비활성화
                    />
                </Stack>
                <Box mt={3}>
                    <Button variant="contained" color="primary" fullWidth>
                        설정 저장
                    </Button>
                </Box>
            </CardContent>
        </Card>
    );
}

export default NetworkConfig;