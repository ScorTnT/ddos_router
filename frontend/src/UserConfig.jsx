import {useState} from 'react';
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

function UserConfig() {
    const [adminId, setAdminId] = useState("root");
    const [adminPassword, setAdminPassword] = useState("");

    return (
        <Card>
            <CardContent>
                <Stack spacing={3}>
                    <Typography variant="h6" gutterBottom>
                        관리자 설정
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

export default UserConfig;