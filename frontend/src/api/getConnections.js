export async function getConnections() {
    try {
        const endpoint = window.location.hostname;
        const response = await fetch(`http://${endpoint}:2024/connections`, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.text();
    } catch (error) {
        console.error('HTTP error! status:', error);
        return "";
    }
}
