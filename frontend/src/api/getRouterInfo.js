const API_URL = import.meta.env.VITE_API_URL

export async function getRouterInfo() {
    try {
        const response = await fetch(`${API_URL}/router_info`);
        if (!response.ok) {
            throw new Error('Failed to fetch router information');
        }
        const routerData = await response.json();
        return routerData;
    } catch (error) {
        console.error('Error fetching router information:', error);
        throw error;
    }
}