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
    }
})
