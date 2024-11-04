export async function attemptLogin(username, password) {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/login`, {
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