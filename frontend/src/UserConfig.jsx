import { useState } from 'react';
import {
    Button,
    Card,
    CardContent,
    TextField,
    Stack,
    Typography,
    Checkbox,
    FormControlLabel
} from "@mui/material";
import { adminConfig, SaveAdminConfig } from './api/adminConfig';

function UserConfig() {
    const [adminId, setAdminId] = useState("root"); // 현재 관리자 계정 (예시: root)
    const [newAdminId, setNewAdminId] = useState("");
    const [newAdminPassword, setNewAdminPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const function (params) {
        
    }
    useEffect(()=>{
        fetchInternetData();
    }, []);
    return (
        <Card>
            <CardContent>
                <Stack spacing={3}>
                    <Typography variant="h6" gutterBottom>
                        관리자 설정
                    </Typography>
                    
                    <Typography variant="subtitle1">
                        현재 관리자 계정: {adminId}
                    </Typography>
                    
                    <TextField
                        label="새 관리자 계정"
                        variant="outlined"
                        value={newAdminId}
                        onChange={(e) => setNewAdminId(e.target.value)}
                        fullWidth
                    />
                    
                    <TextField
                        label="새 관리자 암호"
                        type={showPassword ? "text" : "password"}
                        variant="outlined"
                        value={newAdminPassword}
                        onChange={(e) => setNewAdminPassword(e.target.value)}
                        fullWidth
                    />

                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={showPassword}
                                onChange={() => setShowPassword(!showPassword)}
                                color="primary"
                            />
                        }
                        label="암호보기"
                    />
                    
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

export default UserConfig;
