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
import { getArpInfo } from './api/arpConfig'; 
import { getArpNow } from './api/arpNow';
import { getRouterInfo } from './api/getRouterInfo';
function NetworkConfig() {
    const [connectionType, setConnectionType] = useState("");
    const [ipAddress, setIpAddress] = useState("");
    const [subnetMask, setSubnetMask] = useState("");
    const [gateway, setGateway] = useState("");
    const [primaryDNS, setPrimaryDNS] = useState("");
    const [secondaryDNS, setSecondaryDNS] = useState("");
    const [wanMacAddress, setWanMacAddress] = useState("");
    const [mtu, setMtu] = useState("");
    const [manualDns, setManualDns] = useState();
    const [macAddressChange, setMacAddressChange] = useState();
    const [manualMtu, setManualMtu] = useState();
    const [macOptions, setMacOptions] = useState([]);
    const dhcpLabel = "dhcp 설정 사용중";
    const fetchInternetData = async () => {
        const data = await LoadInternetConfig();
        if (data) {
            setConnectionType(data.connection_type);
            setIpAddress(data.ip_addr || "error");
            setSubnetMask(data.netmask || "error");
            setGateway(data.gateway || "error");
            if (data.connection_type == "dhcp") {
                setIpAddress(dhcpLabel);
                setSubnetMask(dhcpLabel);
                setGateway(dhcpLabel);
            }
            setPrimaryDNS(data.dns_list[0] || "error");
            setSecondaryDNS(data.dns_list[1] || "error");
            setWanMacAddress(data.mac_addr || "");
            setMtu(data.mtu || "error");
            setManualDns(data.is_custom_dns || false);
            setMacAddressChange(data.is_custom_mac || false);
            setManualMtu(data.manualMtu || false);
        }
    };
    const fetchArpData = async () => {
        const routerData = await getRouterInfo();
        let routerMac = null;
        if (routerData) {
            const macRow = routerData.find((row) => row.name === "MAC 주소");
            if (macRow) {
                routerMac = macRow.value;
            }
        }
        const data = await getArpInfo();
        const data1 = await getArpNow();
        console.log(data1);
        if (data) {
            const macAddresses = data.map((row) => row.mac);
            const uniqueMacAddresses = Array.from(new Set(macAddresses));
            const updatedMacOptions = [routerMac, ...uniqueMacAddresses.filter(mac => mac !== routerMac)];
            setMacOptions(updatedMacOptions);
        }

    }
    const saveConfig = async () => {
        const configData = {
            connection_type: connectionType,
            ip_addr: ipAddress,
            netmask: subnetMask,
            gateway: gateway,
            dns_list: [primaryDNS, secondaryDNS],
            mac_addr: wanMacAddress,
            mtu: mtu,
            is_custom_dns: manualDns,
            is_custom_mac: macAddressChange,
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
        fetchArpData();
    }, []);

    // macAddressChange가 비활성화되면 wanMacAddress를 첫 번째 항목으로 설정
    useEffect(() => {
        if (!macAddressChange) {
            setWanMacAddress(macOptions[0]);
        }
    }, [macAddressChange,macOptions]);

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
                            onChange={(e) => {
                                setConnectionType(e.target.value)
                                if (e.target.value === "dhcp") {
                                    setIpAddress(dhcpLabel);
                                    setSubnetMask(dhcpLabel);
                                    setGateway(dhcpLabel);
                                }
                                else if( e.target.value === "static"){
                                    setIpAddress("");
                                    setSubnetMask("");
                                    setGateway("");
                                }
                            }}
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
                            {macOptions.map((mac, index) => (
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