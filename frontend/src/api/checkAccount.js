const API_URL = import.meta.env.VITE_API_URL

export async function attemptLogin(username, password) {
    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({ username, password })
        });

        console.log(await response.json());

        return (response.status === 200);
    } catch {
        return false;
    }
}