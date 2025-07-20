// Simple test to check server connectivity
const baseUrl = 'http://localhost:8080'

console.log('🔍 Testing server connectivity...\n')

async function testServerConnectivity() {
  try {
    console.log('Testing basic server response...')
    
    // Test basic endpoint
    const response = await fetch(`${baseUrl}/notifications`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    console.log(`Response status: ${response.status}`)
    console.log(`Response headers: ${JSON.stringify(Object.fromEntries(response.headers))}`)
    
    const text = await response.text()
    console.log(`Response body: ${text}`)
    
    if (response.ok) {
      console.log('✅ Server is responding correctly')
      try {
        const data = JSON.parse(text)
        console.log('📋 Parsed response:', JSON.stringify(data, null, 2))
      } catch (e) {
        console.log('ℹ️ Response is not JSON')
      }
    } else {
      console.log(`❌ Server responded with error: ${response.status}`)
    }
    
  } catch (error) {
    console.error('❌ Connection failed:', error.message)
    console.error('Stack trace:', error.stack)
  }
}

testServerConnectivity()
