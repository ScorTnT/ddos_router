import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
export default defineConfig({
    base: '/ddos/',
    plugins: [react()],
    server: {
        port: 5173,
        host: '0.0.0.0',
        https: false,
        proxy: {
            '/api': {
                target: 'https://192.168.2.1', // HTTPS로 변경
                changeOrigin: true,
                secure: false, // 자체 서명 인증서 허용
                configure: (proxy) => {
                    proxy.on('proxyReq', (proxyReq) => {
                        // Origin 헤더를 target 주소로 변경
                        proxyReq.setHeader('Origin', 'https://192.168.2.1');
                    });
                }
            }
        }
    }
})
