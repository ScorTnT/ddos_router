import { useEffect } from "react";
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
    TableCell,
    Paper,
    Typography,
    Tooltip,
    IconButton
} from "@mui/material";
function BlackList(){
    const [blockedIP, setBlockedIP] = useState([]);
    useEffect(() => {
        if(blockedIP){
            
        }
    }, []);
    return (
        <Stack spacing={3}>
            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        차단된 IP 목록
                    </Typography>

                    <TableContainer 
                        component={Paper}
                        sx={{
                            overflowY : 'auto',
                        }}
                    >
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 'bold' }}>IP</TableCell>
                                    <TableCell sx={{ fontWeight: 'bold' }}>차단된 시각</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                <TableRow>
                                    <TableCell>     </TableCell>
                                    <TableCell>     </TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
        </Stack>
    )
};

export default BlackList;

