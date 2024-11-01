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

function NetworkConfig() {
    const [stat, setStat] = useState("");
    const [privateIp, setPrivateIp] = useState("");
    const externIp = <></>

    return (
        <Card>
            <CardContent>
                <Stack spacing={3}>
                    <Typography variant="h6" gutterBottom>
                        기본 설정
                    </Typography>

                    <FormControl fullWidth>
                        <InputLabel>무선 채널</InputLabel>
                        <Select
                            value={6}
                            label="무선 채널"
                        >
                            <MenuItem value={1}>채널 1 (2.4GHz)</MenuItem>
                            <MenuItem value={6}>채널 6 (2.4GHz)</MenuItem>
                            <MenuItem value={11}>채널 11 (2.4GHz)</MenuItem>
                        </Select>
                    </FormControl>

                    <FormControl fullWidth>
                        <InputLabel>보안 모드</InputLabel>
                        <Select
                            value="wpa3"
                            label="보안 모드"
                        >
                            <MenuItem value="wpa2">WPA2-PSK</MenuItem>
                            <MenuItem value="wpa3">WPA3-PSK</MenuItem>
                            <MenuItem value="mixed">WPA2/WPA3 혼합</MenuItem>
                        </Select>
                    </FormControl>

                    <FormControlLabel
                        control={<Switch defaultChecked/>}
                        label="게스트 네트워크 활성화"
                    />

                    <FormControlLabel
                        control={<Switch defaultChecked/>}
                        label="QoS 활성화"
                    />

                    <FormControlLabel
                        control={<Switch defaultChecked/>}
                        label="방화벽 활성화"
                    />

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

export default NetworkConfig;