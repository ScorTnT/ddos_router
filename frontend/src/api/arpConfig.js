const API_URL = import.meta.env.VITE_API_URL

export async function getArpInfo() {
    try {
        const response = await fetch(`${API_URL}/arp`);
        if (!response.ok) {
            throw new Error('Failed to fetch arp');
        }
        const arpData = await response.json();
        return arpData;
    } catch (error) {
        console.error('Error fetching arp data:', error);
        throw error;
    }
}