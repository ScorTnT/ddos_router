export async function getConnections() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/connections`);
        //return await response.text();
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('연결 오류:', error);
        return '';
    }
}