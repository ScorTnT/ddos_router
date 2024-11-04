export async function LoadIntranetConfig() {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/api/internet`);

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

export async function SaveIntranetConfig(intranetConfig) {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/api/internet`, {
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