export async function fetchConnectionData() {
    try {
      const endpoint = window.location.hostname;
      const response = await fetch(`http://${endpoint}:2024/connections`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching connection data:', error);
      return [];
    }
  }
  
  export function analyzeConnectionData(data) {
    const protocolCount = {};
    const stateCount = {};
    let totalBytesSent = 0;
    let totalBytesReceived = 0;
  
    data.forEach(connection => {
      // Count protocols
      protocolCount[connection.protocol] = (protocolCount[connection.protocol] || 0) + 1;
  
      // Count states
      stateCount[connection.state] = (stateCount[connection.state] || 0) + 1;
  
      // Sum bytes
      totalBytesSent += parseInt(connection.bytes_sent) || 0;
      totalBytesReceived += parseInt(connection.bytes_received) || 0;
    });
  
    return {
      protocolCount,
      stateCount,
      totalBytesSent,
      totalBytesReceived,
      connectionCount: data.length
    };
  }
  