export async function getHardware() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/hardware`);
        return await response.text();
    } catch (error) {
        console.error('연결 오류:', error);
        return '';
    }
}