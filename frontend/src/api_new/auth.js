import BaseAPI from './base.js';

/**
 * Authentication API class
 */
class AuthAPI extends BaseAPI {
  /**
   * Login to the system
   * @param {string} username - Username
   * @param {string} password - Password
   * @returns {Promise<Object>} - Login response
   */
  async login(username, password) {
    try {
      const response = await this.post('/api/auth/login', {
        username,
        password
      });
      
      // If login successful and response contains session info, store it
      if (response && response.sessionId) {
        this.setSessionId(response.sessionId);
      }
      
      return response;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  /**
   * Logout from the system
   * @returns {Promise<Object>} - Logout response
   */
  async logout() {
    try {
      const response = await this.post('/api/auth/logout');
      
      // Clear session ID after logout
      this.setSessionId(null);
      
      return response;
    } catch (error) {
      console.error('Logout error:', error);
      throw error;
    }
  }

  /**
   * Check if user is authenticated
   * @returns {boolean} - Authentication status
   */
  isAuthenticated() {
    return !!this.getSessionId();
  }
}

export default AuthAPI;