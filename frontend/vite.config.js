import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
export default defineConfig({
    base: '/ddos/',
    plugins: [react()],
    server: {
        port: 5173,
        host: '0.0.0.0',
        https: false, // 임시로 false로 설정
        proxy: {
            '/api': {
                target: 'http://192.168.2.1',
                changeOrigin: true,
                secure: false,
                configure: (proxy) => {
                    proxy.on('proxyReq', (proxyReq) => {
                        // Origin 헤더를 target 주소로 변경
                        proxyReq.setHeader('Origin', 'http://192.168.2.1');
                    });
                }
            }
        }
    }
})
