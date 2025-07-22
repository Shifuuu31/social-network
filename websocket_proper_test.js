const WebSocket = require('ws');

console.log('🔌 Testing WebSocket Notification System (Corrected)');
console.log('====================================================');

// Test WebSocket connection with proper message format
const ws = new WebSocket('ws://localhost:8080/connect?user_id=2');

ws.on('open', function open() {
    console.log('✅ WebSocket connection established as user 2');
    
    // Send notification subscription message in proper Message format
    const subscribeMessage = {
        sender_id: 2,
        receiver_id: 0,
        group_id: 0,
        content: "notification_subscribe",
        type: "notification_subscribe"
    };
    
    ws.send(JSON.stringify(subscribeMessage));
    console.log('📤 Sent notification subscription:', subscribeMessage);
});

ws.on('message', function message(data) {
    try {
        const parsedData = JSON.parse(data);
        console.log('📥 Received message:', JSON.stringify(parsedData, null, 2));
        
        // Check if this is a success response for subscription
        if (parsedData.status === "success" && parsedData.message === "subscribed to notifications") {
            console.log('🎉 Successfully subscribed to notifications!');
            console.log('📡 Now listening for real-time notifications...');
        }
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

// Keep the test running for 15 seconds to listen for real-time notifications
setTimeout(() => {
    console.log('⏰ Test completed - closing WebSocket connection');
    console.log('💡 The WebSocket connection worked properly!');
    ws.close();
}, 15000);

console.log('📡 Listening for real-time notifications for 15 seconds...');
console.log('💡 You can trigger notifications from another terminal using API calls');
