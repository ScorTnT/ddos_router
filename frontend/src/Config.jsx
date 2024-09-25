import React, { useState } from 'react';

function Config() {
  const [settings, setSettings] = useState({
    theme: 'light',
    notifications: true,
    language: 'ko'
  });

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setSettings(prevSettings => ({
      ...prevSettings,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  return (
    <div className="Config">
      <h1>설정 페이지</h1>
      <form>
        <div>
          <label htmlFor="theme">테마:</label>
          <select
            id="theme"
            name="theme"
            value={settings.theme}
            onChange={handleChange}
          >
            <option value="light">라이트</option>
            <option value="dark">다크</option>
          </select>
        </div>
        <div>
          <label htmlFor="notifications">알림:</label>
          <input
            type="checkbox"
            id="notifications"
            name="notifications"
            checked={settings.notifications}
            onChange={handleChange}
          />
        </div>
        <div>
          <label htmlFor="language">언어:</label>
          <select
            id="language"
            name="language"
            value={settings.language}
            onChange={handleChange}
          >
            <option value="ko">한국어</option>
            <option value="en">English</option>
          </select>
        </div>
      </form>
      <div>
        <h2>현재 설정:</h2>
        <p>테마: {settings.theme}</p>
        <p>알림: {settings.notifications ? '켜짐' : '꺼짐'}</p>
        <p>언어: {settings.language === 'ko' ? '한국어' : 'English'}</p>
      </div>
    </div>
  );
}

export default Config;
