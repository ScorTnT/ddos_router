import React, { useState, useEffect } from 'react';
import { fetchConnectionData, analyzeConnectionData } from './statistics';
import { Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

function Config() {
  const [connectionStats, setConnectionStats] = useState(null);

  useEffect(() => {
    async function loadData() {
      const data = await fetchConnectionData();
      const stats = analyzeConnectionData(data);
      setConnectionStats(stats);
    }
    loadData();
  }, []);

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Connection Statistics',
      },
    },
  };

  const protocolChartData = connectionStats ? {
    labels: Object.keys(connectionStats.protocolCount),
    datasets: [
      {
        label: 'Protocol Count',
        data: Object.values(connectionStats.protocolCount),
        backgroundColor: 'rgba(53, 162, 235, 0.5)',
      },
    ],
  } : null;

  const stateChartData = connectionStats ? {
    labels: Object.keys(connectionStats.stateCount),
    datasets: [
      {
        label: 'State Count',
        data: Object.values(connectionStats.stateCount),
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
      },
    ],
  } : null;

  return (
    <div className="Config">
      <h1>Network Connection Statistics</h1>
      {connectionStats ? (
        <div>
          <p>Total Connections: {connectionStats.connectionCount}</p>
          <p>Total Bytes Sent: {connectionStats.totalBytesSent}</p>
          <p>Total Bytes Received: {connectionStats.totalBytesReceived}</p>
          
          <div style={{width: '500px', height: '300px'}}>
            <h2>Protocol Distribution</h2>
            <Bar options={chartOptions} data={protocolChartData} />
          </div>
          
          <div style={{width: '500px', height: '300px'}}>
            <h2>Connection State Distribution</h2>
            <Bar options={chartOptions} data={stateChartData} />
          </div>
        </div>
      ) : (
        <p>Loading connection data...</p>
      )}
    </div>
  );
}

export default Config;
