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
<<<<<<< HEAD
import apiClient from './api.js';
=======
import api from './api.js';
>>>>>>> abbd53f9fddff4e463e9be79891f7bc96cb18a04

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
        setError('');
        setIsLoading(true);

        if (!formData.username.trim() || !formData.password.trim()) {
            setError('아이디와 비밀번호를 입력해주세요.');
            setIsLoading(false);
            return;
        }

        try {
<<<<<<< HEAD
            console.log('[DEBUG] Submitting login with:', formData);
            const response = await apiClient.login(formData.username, formData.password);
            console.log('Login successful:', response);
=======
            await api.login(formData.username, formData.password);
>>>>>>> abbd53f9fddff4e463e9be79891f7bc96cb18a04
            setIsLoggedIn(true);
        } catch (err) {
            setError(err.message || '아이디 또는 비밀번호가 올바르지 않습니다.');
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