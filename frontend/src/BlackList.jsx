import { useEffect, useState } from "react";
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
    IconButton,
    TableRow
} from "@mui/material";
import { getBlockedIP, blockIP, unblockIP } from './api/getBlockedIP';
function BlackList(){
    const [blockedIP, setBlockedIP] = useState([]);
    const [blackList, setBlackList] = useState([]);
    const [whiteList, setWhiteList] = useState([]);
    useEffect(() => {
        // if(){

        // }
    }, []);
    return <>
        <Stack spacing={3}>
            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        차단 IP 목록 (Dropped IP)
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
                                    <TableCell> api 답변에 따라 수정 예정 </TableCell>
                                    <TableCell> api 답변에 따라 수정 예정 </TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>

            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        영구 차단 IP 목록 (IP BlackList)
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
                                    <TableCell sx={{ fontWeight: 'bold' }}>추가된 시각</TableCell>
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

            <Card>
                <CardContent>
                    <Typography variant="h6" gutterBottom>
                        허용 IP 목록 (IP WhiteList)
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
                                    <TableCell sx={{ fontWeight: 'bold' }}>허용된 시각</TableCell>
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
    </>
}

export default BlackList;