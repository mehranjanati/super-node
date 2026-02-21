const WebSocket = require('ws');
const ws = new WebSocket('ws://localhost:3000/ws/market');

ws.on('open', function open() {
  console.log('connected');
});

ws.on('message', function incoming(data) {
  console.log('received: %s', data);
});

ws.on('error', function error(err) {
  console.error('error:', err);
});

setTimeout(() => {
    console.log('Closing connection...');
    ws.close();
}, 10000); // Listen for 10 seconds
