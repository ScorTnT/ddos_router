// src/api/getRouterInfo.js
export async function getRouterInfo() {
    try {
        const response = await fetch('http://greatmandu.asuscomm.com:2024/router_info');
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