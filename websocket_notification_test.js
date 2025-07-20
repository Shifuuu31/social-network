const WebSocket = require('ws');

console.log('🔌 Testing WebSocket Notification System');
console.log('=====================================');

// Test WebSocket connection with proper notification subscription
const ws = new WebSocket('ws://localhost:8080/connect?user_id=2');

ws.on('open', function open() {
    console.log('✅ WebSocket connection established as user 2');
    
    // Subscribe to notifications using the correct format
    const subscribeMessage = {
        type: 'notification_subscribe'
    };
    
    ws.send(JSON.stringify(subscribeMessage));
    console.log('📤 Sent notification subscription:', subscribeMessage);
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
    process.exit(0);
});

// Keep the test running for 10 seconds to listen for real-time notifications
setTimeout(() => {
    console.log('⏰ Test completed - closing WebSocket connection');
    ws.close();
}, 10000);

console.log('📡 Listening for real-time notifications for 10 seconds...');
console.log('💡 You can trigger notifications from another terminal using API calls');
