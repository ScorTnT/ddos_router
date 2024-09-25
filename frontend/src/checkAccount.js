export async function attemptLogin(username, password) {
  try {
    const endpoint = window.location.hostname;
    const response = await fetch(`http://${endpoint}:2024/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams({
        'username': username,
        'password': password
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.text();
    const isLoginSuccessful = result === 'true';

    if (isLoginSuccessful) {
      console.log('로그인 성공');
    } else {
      console.log('로그인 실패');
    }

    return isLoginSuccessful;
  } catch (error) {
    console.error('로그인 시도 중 오류 발생:', error);
    return false;
  }
}
