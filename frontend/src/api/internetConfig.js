const API_URL = import.meta.env.VITE_API_URL

export async function LoadInternetConfig() {
    try {
        const response = await fetch(`${API_URL}/internet`);

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
        const response = await fetch(`${API_URL}/internet`, {
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