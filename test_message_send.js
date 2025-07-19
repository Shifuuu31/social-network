// Test script to send a message via HTTP API
const fetch = require('node-fetch');

async function testSendMessage() {
  console.log('🔍 Testing message sending via HTTP API...');
  
  try {
    // First, let's try to send a message
    const response = await fetch('http://localhost:8080/api/chat/send', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': 'session=your-session-cookie-here' // You'll need to get this from browser
      },
      body: JSON.stringify({
        receiver_id: 2,
        content: 'Test message from Node.js script'
      })
    });

    console.log('Response status:', response.status);
    console.log('Response headers:', Object.fromEntries(response.headers.entries()));
    
    const responseText = await response.text();
    console.log('Raw response:', responseText);
    
    if (response.ok) {
      const data = JSON.parse(responseText);
      console.log('✅ Message sent successfully:', data);
    } else {
      console.log('❌ Failed to send message');
    }
  } catch (error) {
    console.error('❌ Error:', error);
  }
}

testSendMessage(); 