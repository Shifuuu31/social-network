const WebSocket = require('ws');

console.log('🧪 Testing WebSocket Connection');
console.log('===============================');

const ws = new WebSocket('ws://localhost:8080/connect');

ws.on('open', function open() {
    console.log('✅ WebSocket connection established!');
    console.log('📤 Sending test message...');
    
    // Send a test message
    ws.send(JSON.stringify({
        type: 'test',
        message: 'Hello WebSocket!'
    }));
});

ws.on('message', function message(data) {
    console.log('📨 Received message:', data.toString());
});

ws.on('error', function error(err) {
    console.log('❌ WebSocket error:', err.message);
});

ws.on('close', function close(code, reason) {
    console.log('🔌 WebSocket connection closed');
    console.log('Code:', code);
    console.log('Reason:', reason.toString());
});

// Close connection after 5 seconds
setTimeout(() => {
    if (ws.readyState === WebSocket.OPEN) {
        console.log('🔚 Closing WebSocket connection...');
        ws.close();
    }
    process.exit(0);
}, 5000);
