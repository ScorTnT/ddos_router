export async function adminId() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/api/user`);

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