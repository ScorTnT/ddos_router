import BaseAPI from './base.js';

/**
 * Protection API class for security features
 */
class ProtectionAPI extends BaseAPI {
  /**
   * Get protection list
   * @returns {Promise<Object>} - Protection list data
   */
  async getProtectionList() {
    try {
      return await this.get('/api/protection');
    } catch (error) {
      console.error('Get protection list error:', error);
      throw error;
    }
  }

  /**
   * Block an IP address
   * @param {string} ip - IP address to block
   * @returns {Promise<Object>} - Block operation response
   */
  async blockIP(ip) {
    try {
      if (!ip) {
        throw new Error('IP address is required');
      }

      return await this.post('/api/protection/ip/block', {}, { ip });
    } catch (error) {
      console.error('Block IP error:', error);
      throw error;
    }
  }

  /**
   * Unblock an IP address
   * @param {string} ip - IP address to unblock
   * @returns {Promise<Object>} - Unblock operation response
   */
  async unblockIP(ip) {
    try {
      if (!ip) {
        throw new Error('IP address is required');
      }

      return await this.post('/api/protection/ip/unblock', {}, { ip });
    } catch (error) {
      console.error('Unblock IP error:', error);
      throw error;
    }
  }

  /**
   * Toggle IP block status
   * @param {string} ip - IP address to toggle
   * @param {boolean} shouldBlock - True to block, false to unblock
   * @returns {Promise<Object>} - Toggle operation response
   */
  async toggleIPBlock(ip, shouldBlock) {
    try {
      if (shouldBlock) {
        return await this.blockIP(ip);
      } else {
        return await this.unblockIP(ip);
      }
    } catch (error) {
      console.error('Toggle IP block error:', error);
      throw error;
    }
  }
}

export default ProtectionAPI;
