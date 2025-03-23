export async function getConnections() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/connections`);
        
        if (!response.ok) {
            throw new Error("Failed to fetch router connection informations");
        }
        
        const connections = await response.json();
        console.log(connections);

        return connections;
    } catch (error) {
        console.error('연결 오류:', error);
        return '';
    }
}