import AuthAPI from './auth.js';
import InformationAPI from './information.js';
import ConfigAPI from './config.js';
import ProtectionAPI from './protection.js';

// new
import Cookies from 'js-cookie';
// npm install js-cookie >> needed

/**
 * Main API client that combines all API modules
 */
class RouterAPIClient {
  constructor(baseURL = '') {
    this.baseURL = baseURL;
    
    // Initialize all API modules
    this.auth = new AuthAPI(baseURL);
    this.information = new InformationAPI(baseURL);
    this.config = new ConfigAPI(baseURL);
    this.protection = new ProtectionAPI(baseURL);
    
    // Keep track of session across all modules
    this._sessionId = null;
  }

  // /**
  //  * Set session ID for all API modules
  //  * @param {string} sessionId - The session ID
  //  */
  // setSessionId(sessionId) {
  //   this._sessionId = sessionId;
  //   this.auth.setSessionId(sessionId);
  //   this.information.setSessionId(sessionId);
  //   this.config.setSessionId(sessionId);
  //   this.protection.setSessionId(sessionId);
  // }

  // /**
  //  * Get current session ID
  //  * @returns {string|null} - The current session ID
  //  */
  // getSessionId() {
  //   return this._sessionId;
  // }

  /**
   * Set session ID for authenticated requests
   * @param {string} sessionId - The session ID
  */
  setSessionId(sessionId) {
    if (sessionId) {
      Cookies.set('sessionId', sessionId, { expires: 1 });
    } else {
      Cookies.remove('sessionId');
    }
  }

  /**
   * Get session ID
   * @returns {string|null} - The current session ID
   */
  getSessionId() {
    return Cookies.get('sessionId') || null;
  }

  /**
   * Login and set session for all modules
   * @param {string} username - Username
   * @param {string} password - Password
   * @returns {Promise<Object>} - Login response
   */
  async login(username, password) {
    try {
      const response = await this.auth.login(username, password);
      
      // Sync session ID across all modules
      const sessionId = this.auth.getSessionId();
      if (sessionId) {
        this.setSessionId(sessionId);
      }
      
      return response;
    } catch (error) {
      console.error('Client login error:', error);
      throw error;
    }
  }

  /**
   * Logout and clear session from all modules
   * @returns {Promise<Object>} - Logout response
   */
  async logout() {
    try {
      const response = await this.auth.logout();
      
      // Clear session from all modules
      this.setSessionId(null);
      
      return response;
    } catch (error) {
      console.error('Client logout error:', error);
      throw error;
    }
  }

  /**
   * Check if user is authenticated
   * @returns {boolean} - Authentication status
   */
  isAuthenticated() {
    return this.auth.isAuthenticated();
  }

  /**
   * Set base URL for all API modules
   * @param {string} baseURL - The base URL
   */
  setBaseURL(baseURL) {
    this.baseURL = baseURL;
    this.auth.baseURL = baseURL;
    this.information.baseURL = baseURL;
    this.config.baseURL = baseURL;
    this.protection.baseURL = baseURL;
  }
}

// Create and export a default instance
const apiClient = new RouterAPIClient();

export default apiClient;
export { RouterAPIClient };
