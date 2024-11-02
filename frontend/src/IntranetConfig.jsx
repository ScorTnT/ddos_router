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
    const intranetStatus = <>intranet</>

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