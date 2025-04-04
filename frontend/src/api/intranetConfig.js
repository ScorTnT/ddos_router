const API_URL = import.meta.env.VITE_API_URL

export async function LoadIntranetConfig() {
    try {
        const response = await fetch(`${API_URL}/intranet`);

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
        const response = await fetch(`${API_URL}/intranet`, {
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