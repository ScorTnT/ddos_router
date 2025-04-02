const API_URL = import.meta.env.VITE_API_URL

export async function getArpNow() {
    try {
        const response = await fetch(`${API_URL}/arp_now`);
        if (!response.ok) {
            throw new Error('Failed to fetch arp_now');
        }
        const arpData = await response.json();
        return arpData;
    } catch (error) {
        console.error('Error fetching arp data:', error);
        throw error;
    }
}