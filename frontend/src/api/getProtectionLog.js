const API_URL = import.meta.env.VITE_API_URL;

// Protectioni Log 조회
export async function getProtectionLogs() {
    try {
        const response = await fetch(`${API_URL}/protection`);
        if (!response.ok) {
            throw new Error("Failed to fetch protection log");
        }
        const protectionLog = await response.json();
        return protectionLog;
    } catch (error) {
        console.error("Error fetching protection logs:", error);
        throw error;
    }
}

// IP 수동 차단
export async function blockIP(ip) {
    try {
        const response = await fetch(`${API_URL}/protection/block`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ ip }),
        });

        if (!response.ok) {
            throw new Error("Failed to block IP");
        }

        const result = await response.json();
        return result;
    } catch (error) {
        console.error("Error blocking IP:", error);
        throw error;
    }
}

// IP 수동 차단 해제
export async function unblockIP(ip) {
    try {
        const response = await fetch(`${API_URL}/protection/unblock`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ ip }),
        });

        if (!response.ok) {
            throw new Error("Failed to unblock IP");
        }

        const result = await response.json();
        return result;
    } catch (error) {
        console.error("Error unblocking IP:", error);
        throw error;
    }
}