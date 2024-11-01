export async function getConnections() {
    try {
        const response = await fetch(`${window.location.origin}:2024/connections`);
        return await response.text();
    } catch (error) {
        console.error('연결 오류:', error);
        return '';
    }
}