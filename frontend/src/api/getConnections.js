export async function getConnections() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/connections`);
        console.log(response);
        return await response.text();
    } catch (error) {
        console.error('연결 오류:', error);
        return '';
    }
}