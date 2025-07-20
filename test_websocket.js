const WebSocket = require('ws');

console.log('🔌 Testing WebSocket Notification System');
console.log('=====================================');

// Test WebSocket connection
const ws = new WebSocket('ws://localhost:8080/connect');

ws.on('open', function open() {
    console.log('✅ WebSocket connection established');
    
    // Send a test message to identify the user
    const testMessage = {
        type: 'user_connect',
        user_id: 2
    };
    
    ws.send(JSON.stringify(testMessage));
    console.log('📤 Sent user connect message:', testMessage);
});

ws.on('message', function message(data) {
    try {
        const parsedData = JSON.parse(data);
        console.log('📥 Received message:', JSON.stringify(parsedData, null, 2));
    } catch (e) {
        console.log('📥 Received raw message:', data.toString());
    }
});

ws.on('error', function error(err) {
    console.error('❌ WebSocket error:', err.message);
});

ws.on('close', function close() {
    console.log('🔌 WebSocket connection closed');
});

// Keep the test running for 10 seconds
setTimeout(() => {
    console.log('⏰ Closing WebSocket connection after test');
    ws.close();
    process.exit(0);
}, 10000);
