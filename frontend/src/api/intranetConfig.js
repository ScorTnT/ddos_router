export async function LoadInternetConfig() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/api/intranet`);

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

export async function SaveInternetConfig(intranetConfig) {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/api/intranet`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(intranetConfig)
        });

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