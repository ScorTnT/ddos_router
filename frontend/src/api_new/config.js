import BaseAPI from './base.js';

/**
 * Configuration API class for router settings
 */
class ConfigAPI extends BaseAPI {
  /**
   * Get LAN configuration
   * @returns {Promise<Object>} - LAN configuration data
   */
  async getLANConfig() {
    try {
      return await this.get('/api/config/lan');
    } catch (error) {
      console.error('Get LAN config error:', error);
      throw error;
    }
  }

  /**
   * Get WAN configuration
   * @returns {Promise<Object>} - WAN configuration data
   */
  async getWANConfig() {
    try {
      return await this.get('/api/config/wan');
    } catch (error) {
      console.error('Get WAN config error:', error);
      throw error;
    }
  }

  /**
   * Get all network configurations
   * @returns {Promise<Object>} - Object containing both LAN and WAN configs
   */
  async getAllConfigs() {
    try {
      const [lanConfig, wanConfig] = await Promise.all([
        this.getLANConfig(),
        this.getWANConfig()
      ]);
      
      return {
        lan: lanConfig,
        wan: wanConfig
      };
    } catch (error) {
      console.error('Get all configs error:', error);
      throw error;
    }
  }
}

export default ConfigAPI;
