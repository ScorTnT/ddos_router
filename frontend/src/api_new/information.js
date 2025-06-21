import BaseAPI from './base.js';

/**
 * Information API class for router data
 */
class InformationAPI extends BaseAPI {
  /**
   * Get router information
   * @returns {Promise<Object>} - Router information
   */
  async getRouterInfo() {
    try {
      return await this.get('/api/information');
    } catch (error) {
      console.error('Get router info error:', error);
      throw error;
    }
  }

  /**
   * Get router ARP table (neighbors)
   * @returns {Promise<Object>} - ARP table data
   */
  async getNeighbors() {
    try {
      return await this.get('/api/information/neighbors');
    } catch (error) {
      console.error('Get neighbors error:', error);
      throw error;
    }
  }

  /**
   * Get router connection information
   * @returns {Promise<Object>} - Connection information
   */
  async getConnections() {
    try {
      return await this.get('/api/information/connections');
    } catch (error) {
      console.error('Get connections error:', error);
      throw error;
    }
  }
}

export default InformationAPI;
