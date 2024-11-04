export async function adminConfig() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/user`);

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
        const response = await fetch(`http://${window.location.hostname}:2024/user`);

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