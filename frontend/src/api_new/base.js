/**
 * Base API class for handling HTTP requests
 */
class BaseAPI {
  constructor(baseURL = '') {
    this.baseURL = baseURL;
    this.sessionId = null;
  }

  /**
   * Set session ID for authenticated requests
   * @param {string} sessionId - The session ID
   */
  setSessionId(sessionId) {
    this.sessionId = sessionId;
  }

  /**
   * Get session ID
   * @returns {string|null} - The current session ID
   */
  getSessionId() {
    return this.sessionId;
  }

  /**
   * Get default headers
   * @returns {Object} - Default headers object
   */
  getHeaders() {
    const headers = {
      'Content-Type': 'application/json',
    };

    if (this.sessionId) {
      headers['X-Session-ID'] = this.sessionId;
    }

    return headers;
  }

  /**
   * Make HTTP request
   * @param {string} endpoint - API endpoint
   * @param {Object} options - Request options
   * @returns {Promise<Object>} - Response data
   */
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: this.getHeaders(),
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        throw new Error(`HTTP Error: ${response.status} ${response.statusText}`);
      }

      // Check if response has content
      const contentType = response.headers.get('Content-Type');
      if (contentType && contentType.includes('application/json')) {
        return await response.json();
      }
      
      return await response.text();
    } catch (error) {
      console.error('API Request Error:', error);
      throw error;
    }
  }

  /**
   * Make GET request
   * @param {string} endpoint - API endpoint
   * @param {Object} params - Query parameters
   * @returns {Promise<Object>} - Response data
   */
  async get(endpoint, params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const url = queryString ? `${endpoint}?${queryString}` : endpoint;
    
    return this.request(url, {
      method: 'GET',
    });
  }

  /**
   * Make POST request
   * @param {string} endpoint - API endpoint
   * @param {Object} data - Request body data
   * @param {Object} params - Query parameters
   * @returns {Promise<Object>} - Response data
   */
  async post(endpoint, data = {}, params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const url = queryString ? `${endpoint}?${queryString}` : endpoint;
    
    return this.request(url, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
}

export default BaseAPI;
