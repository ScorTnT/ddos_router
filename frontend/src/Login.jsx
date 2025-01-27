import {useState} from 'react';
import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    CircularProgress,
    IconButton,
    InputAdornment,
    TextField,
    Typography
} from '@mui/material';
import {Lock, Person, Visibility, VisibilityOff} from '@mui/icons-material';
import PropTypes from 'prop-types';
import { attemptLogin } from './api/checkAccount';

const Login = ({setIsLoggedIn}) => {
    const [formData, setFormData] = useState({
        username: '',
        password: ''
    });
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleChange = (e) => {
        const {name, value} = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
        setError('');
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            const loginSuccess = await attemptLogin(formData.username, formData.password); // 이 둘을 절대 조심
            //const loginSuccess = true;                                                       // 위는 실제 서버에서 로그인 확인, 아래는 무조건 통과
            if (loginSuccess) {
                setIsLoggedIn(true);
            } else {
                setError('아이디 또는 비밀번호가 올바르지 않습니다.');
            }
        } catch (err) {
            setError('로그인 중 오류가 발생했습니다. 다시 시도해주세요.');
            console.error('Login error:', err);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: '#f5f5f5',
                padding: 2
            }}
        >
            <Card sx={{maxWidth: 400, width: '100%'}}>
                <CardContent sx={{padding: 4}}>
                    <Box component="form" onSubmit={handleSubmit}
                         sx={{display: 'flex', flexDirection: 'column', gap: 3}}>
                        <Box sx={{textAlign: 'center', mb: 2}}>
                            <Typography variant="h4" component="h1" gutterBottom>
                                라우터 관리 시스템
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                                계정 정보를 입력하여 로그인하세요
                            </Typography>
                        </Box>

                        <TextField
                            fullWidth
                            label="아이디"
                            name="username"
                            value={formData.username}
                            onChange={handleChange}
                            required
                            slotProps={{
                                input: {
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <Person/>
                                        </InputAdornment>
                                    ),
                                }
                            }}
                        />

                        <TextField
                            fullWidth
                            label="비밀번호"
                            name="password"
                            type={showPassword ? 'text' : 'password'}
                            value={formData.password}
                            onChange={handleChange}
                            required
                            slotProps={{
                                input: {
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <Lock/>
                                        </InputAdornment>
                                    ),
                                    endAdornment: (
                                        <InputAdornment position="end">
                                            <IconButton
                                                onClick={() => setShowPassword(!showPassword)}
                                                edge="end"
                                            >
                                                {showPassword ? <VisibilityOff/> : <Visibility/>}
                                            </IconButton>
                                        </InputAdornment>
                                    ),
                                }
                            }}
                        />

                        {error && (
                            <Alert severity="error" sx={{width: '100%'}}>
                                {error}
                            </Alert>
                        )}

                        <Button
                            type="submit"
                            variant="contained"
                            size="large"
                            disabled={isLoading}
                            sx={{
                                mt: 2,
                                height: 48,
                                backgroundColor: '#1976d2',
                                '&:hover': {
                                    backgroundColor: '#1565c0'
                                }
                            }}
                        >
                            {isLoading ? (
                                <CircularProgress size={24} color="inherit"/>
                            ) : (
                                '로그인'
                            )}
                        </Button>
                    </Box>
                </CardContent>
            </Card>
        </Box>
    );
};

Login.propTypes = {
    setIsLoggedIn: PropTypes.func.isRequired
};

export default Login;