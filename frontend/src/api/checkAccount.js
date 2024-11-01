export async function attemptLogin(username, password) {
    try {
        const response = await fetch(`http://${window.location.hostname}:2024/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({ username, password })
        });

        const result = await response.text();
        const isLoginSuccessful = result === 'true';

        console.log(isLoginSuccessful ? '로그인 성공' : '로그인 실패');
        return isLoginSuccessful;

    } catch {
        return false;
    }
}