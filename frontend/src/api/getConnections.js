const API_URL = import.meta.env.VITE_API_URL

export async function getConnections() {
    try {
        const response = await fetch(`${API_URL}//connections`);
        
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