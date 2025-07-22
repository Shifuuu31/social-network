// Test WebSocket notifications functionality
const WebSocket = require('ws')

async function testWebSocketNotifications() {
  console.log('🌐 Testing WebSocket Notifications\n')
  
  return new Promise((resolve, reject) => {
    const ws = new WebSocket('ws://localhost:8080/connect')
    let connected = false
    
    // Set timeout for connection
    const timeout = setTimeout(() => {
      if (!connected) {
        console.log('❌ WebSocket connection timeout')
        ws.close()
        resolve()
      }
    }, 5000)
    
    ws.on('open', () => {
      connected = true
      clearTimeout(timeout)
      console.log('✅ WebSocket connected successfully')
      
      // Test sending a message
      ws.send(JSON.stringify({
        type: 'ping',
        data: 'test'
      }))
      
      // Close after a short delay
      setTimeout(() => {
        console.log('🔌 Closing WebSocket connection')
        ws.close()
        resolve()
      }, 2000)
    })
    
    ws.on('message', (data) => {
      try {
        const message = JSON.parse(data.toString())
        console.log('📩 Received message:', message)
      } catch (error) {
        console.log('📩 Received raw data:', data.toString())
      }
    })
    
    ws.on('error', (error) => {
      console.log('❌ WebSocket error:', error.message)
      clearTimeout(timeout)
      resolve()
    })
    
    ws.on('close', () => {
      console.log('🔌 WebSocket connection closed')
      clearTimeout(timeout)
      resolve()
    })
  })
}

// Run the test
testWebSocketNotifications()
  .then(() => console.log('\n🎉 WebSocket test completed!'))
  .catch(console.error)
