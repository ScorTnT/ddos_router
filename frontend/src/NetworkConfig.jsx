import { useState, useEffect } from 'react';
import {
    Box,
    Button,
    Card,
    CardContent,
    FormControl,
    FormControlLabel,
    InputLabel,
    MenuItem,
    Radio,
    RadioGroup,
    Select,
    Stack,
    Switch,
    TextField,
    Typography
} from "@mui/material";
import { LoadInternetConfig, SaveInternetConfig } from './api/internetConfig';

function NetworkConfig() {
    const [connectionType, setConnectionType] = useState("dhcp"); // 동적 IP 방식 기본값
    const [ipAddress, setIpAddress] = useState("220.66.87.40");
    const [subnetMask, setSubnetMask] = useState("255.255.255.0");
    const [gateway, setGateway] = useState("220.66.87.2");
    const [primaryDNS, setPrimaryDNS] = useState("8.8.8.8");
    const [secondaryDNS, setSecondaryDNS] = useState("8.8.4.4");
    const [wanMacAddress, setWanMacAddress] = useState(""); // 선택한 MAC 주소 저장
    const [mtu, setMtu] = useState("1500");
    const [manualDns, setManualDns] = useState(true);
    const [macAddressChange, setMacAddressChange] = useState(true);
    const [manualMtu, setManualMtu] = useState(false);
    const fetchInternetData = async () => {
        const data = await LoadInternetConfig();
        if (data) {
            setConnectionType(data.connection_type || "error");
            setIpAddress(data.ip_addr || "error");
            setSubnetMask(data.netmask || "error");
            setGateway(data.gateway || "error");
            setPrimaryDNS(data.dns_list[0] || "error");
            setSecondaryDNS(data.dns_list[1] || "error");
            setWanMacAddress(data.mac_addr || "");
            setMtu(data.mtu || "error");
            setManualDns(data.use_custom_dns || false);
            setMacAddressChange(data.clone_mac || false);
            setManualMtu(data.manualMtu || false);
        }
    };
    const saveConfig = async () => {
        const configData = {
            connection_type: connectionType,
            ip_addr: ipAddress,
            netmask: subnetMask,
            gateway: gateway,
            dns_list: [primaryDNS, secondaryDNS],
            mac_addr: wanMacAddress,
            mtu: mtu,
            use_custom_dns: manualDns,
            clone_mac: macAddressChange,
            manualMtu: manualMtu,
        };
    
        const result = await SaveInternetConfig(configData);
        if (result) {
            alert("설정이 성공적으로 저장되었습니다!");
        } else {
            alert("설정 저장에 실패했습니다.");
        }
    };
    useEffect(()=>{
        fetchInternetData();
    }, []);

    // 예시 MAC 주소 리스트
    const macAddressOptions = [
        "FF:AA:AA:99:6F:TE",
        "C4:C3:KK:34:55:AR",
        "TT:BE:JE:12:34:5K"
    ];

    // macAddressChange가 비활성화되면 wanMacAddress를 첫 번째 항목으로 설정
    useEffect(() => {
        if (!macAddressChange) {
            setWanMacAddress(macAddressOptions[0]);
        }
    }, [macAddressChange]);

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
                                value="dhcp"
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
                        disabled={connectionType === "dhcp"} // 동적 IP 방식일 때 비활성화
                    />
                    <TextField
                        label="서브넷 마스크"
                        value={subnetMask}
                        onChange={(e) => setSubnetMask(e.target.value)}
                        fullWidth
                        disabled={connectionType === "dhcp"} // 동적 IP 방식일 때 비활성화
                    />
                    <TextField
                        label="기본 게이트웨이"
                        value={gateway}
                        onChange={(e) => setGateway(e.target.value)}
                        fullWidth
                        disabled={connectionType === "dhcp"} // 동적 IP 방식일 때 비활성화
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
                    <FormControl fullWidth disabled={!macAddressChange}>
                        <InputLabel>WAN MAC 주소</InputLabel>
                        <Select
                            value={wanMacAddress}
                            onChange={(e) => setWanMacAddress(e.target.value)}
                        >
                            {macAddressOptions.map((mac, index) => (
                                <MenuItem key={index} value={mac}>
                                    {mac}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>

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
                    <Button variant="contained" color="primary" fullWidth onClick={saveConfig}>
                        설정 저장
                    </Button>
                </Box>
            </CardContent>
        </Card>
    );
}

export default NetworkConfig;