const API_URL = import.meta.env.VITE_API_URL

export async function adminConfig() {
    try {
        const response = await fetch(`${API_URL}/user`);

        if (response.ok) {
            return await response.json();
        } else {
            return null;
        }
    } catch (error) {
        console.error(error);
        return null;
    }
}export async function SaveAdminConfig() {
    try {
        const response = await fetch(`${API_URL}/user`);

        if (response.ok) {
            return await response.json();
        } else {
            return null;
        }
    } catch (error) {
        console.error(error);
        return null;
    }
}