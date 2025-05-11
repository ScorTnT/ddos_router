const API_URL = import.meta.env.VITE_API_URL

export async function getHardware() {
    try {
        const response = await fetch(`${API_URL}/hardware`);
        return await response.text();
    } catch (error) {
        console.error('연결 오류:', error);
        return '';
    }
}