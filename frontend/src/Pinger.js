import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { Table } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';

const PingTable = () => {
  const [data, setData] = useState([]);
  const [columnWidths, setColumnWidths] = useState({
    status: 70,
    id_container: 200,
    name: 200,
    ip: 50,
    ping: 50,
    last_success: 95,
  });

  const apiUrl = process.env.REACT_APP_API_URL;

  const fetchData = useCallback(async () => {
    try {
      const response = await axios.get(apiUrl);
      setData(response.data.data || []);
    } catch (error) {
      console.error('Error fetching data:', error);
      setData([]);
    }
  }, [apiUrl]);
  
  useEffect(() => {
    fetchData(); // Вызываем один раз при монтировании
    const intervalId = setInterval(fetchData, 10000); // Каждые 10 секунд
    return () => clearInterval(intervalId);
  }, [fetchData]); // Теперь зависимость `fetchData` не меняется
  
  const handleMouseDown = (e, column) => {
    const startX = e.clientX;
    const startWidth = columnWidths[column];
    e.preventDefault();

    const onMouseMove = (moveEvent) => {
      const newWidth = startWidth + moveEvent.clientX - startX;
      setColumnWidths((prevWidths) => ({
        ...prevWidths,
        [column]: Math.max(newWidth, 30),
      }));
    };

    const onMouseUp = () => {
      document.removeEventListener('mousemove', onMouseMove);
      document.removeEventListener('mouseup', onMouseUp);
    };

    document.addEventListener('mousemove', onMouseMove);
    document.addEventListener('mouseup', onMouseUp);
  };

  return (
    <div className="container mt-4">
      <h2>Ping Data</h2>
      <Table striped bordered hover style={{ tableLayout: 'auto', width: '100%' }}>
        <thead>
          <tr>
            <th
              style={{ width: columnWidths.status }}
              onMouseDown={(e) => handleMouseDown(e, 'status')}
            >
              Status
            </th>
            <th
              style={{ width: columnWidths.id_container }}
              onMouseDown={(e) => handleMouseDown(e, 'id_container')}
            >
              Container ID
            </th>
            <th
              style={{ width: columnWidths.name }}
              onMouseDown={(e) => handleMouseDown(e, 'name')}
            >
              Container Name
            </th>
            <th
              style={{ width: columnWidths.ip }}
              onMouseDown={(e) => handleMouseDown(e, 'ip')}
            >
              IP Address
            </th>
            <th
              style={{ width: columnWidths.ping }}
              onMouseDown={(e) => handleMouseDown(e, 'ping')}
            >
              Ping Time
            </th>
            <th
              style={{ width: columnWidths.last_success }}
              onMouseDown={(e) => handleMouseDown(e, 'last_success')}
            >
              Last Success
            </th>
          </tr>
        </thead>
        <tbody>
          {data.length > 0 ? (
            data.map((item) => (
              <tr key={item.id_container}>
                <td style={{ maxWidth: columnWidths.status, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {item.status ? 'Active' : 'Inactive'}
                </td>
                <td style={{ maxWidth: columnWidths.id_container, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {item.id_container}
                </td>
                <td style={{ maxWidth: columnWidths.name, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {item.name}
                </td>
                <td style={{ maxWidth: columnWidths.ip, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {item.ip || 'N/A'}
                </td>
                <td style={{ maxWidth: columnWidths.ping, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {item.response_time.toFixed(3)} ms
                </td>
                <td style={{ maxWidth: columnWidths.last_success, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {item.last_success}
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="6" className="text-center">No data available</td>
            </tr>
          )}
        </tbody>
      </Table>
    </div>
  );
};

export default PingTable;
