// Test script to verify groups functionality with notifications
const baseUrl = 'http://localhost:8080'

async function testGroupsWithNotifications() {
  console.log('🧪 Testing Groups Functionality with Notifications\n')
  
  // Test 1: Fetch groups
  try {
    console.log('📋 Testing groups fetch...')
    const response = await fetch(`${baseUrl}/groups/group/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: "1",
        start: -1,
        n_items: 20,
        type: "user"
      })
    })
    
    if (response.ok) {
      const groups = await response.json()
      console.log(`✅ Successfully fetched ${groups?.length || 0} groups`)
      if (groups?.length > 0) {
        console.log('   Sample group:', groups[0].title)
      }
    } else {
      console.log('❌ Failed to fetch groups:', response.status)
    }
  } catch (error) {
    console.log('❌ Error fetching groups:', error.message)
  }
  
  // Test 2: Create a test group
  try {
    console.log('\n🏗️ Testing group creation...')
    const response = await fetch(`${baseUrl}/groups/group/new`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: "Test Notifications Group",
        description: "Testing group with notifications",
        user_id: "1"
      })
    })
    
    if (response.ok) {
      const result = await response.json()
      console.log('✅ Successfully created test group')
      console.log('   Group ID:', result.group_id || 'N/A')
    } else {
      console.log('❌ Failed to create group:', response.status)
    }
  } catch (error) {
    console.log('❌ Error creating group:', error.message)
  }
  
  // Test 3: Test notification count
  try {
    console.log('\n🔔 Testing notification count...')
    const response = await fetch(`${baseUrl}/notifications/unread-count?user_id=1`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    })
    
    if (response.ok) {
      const result = await response.json()
      console.log('✅ Successfully fetched notification count')
      console.log('   Unread count:', result.count || 0)
    } else {
      console.log('❌ Failed to fetch notification count:', response.status)
    }
  } catch (error) {
    console.log('❌ Error fetching notification count:', error.message)
  }
  
  console.log('\n🎉 Groups and notifications test completed!')
}

// Run the test
testGroupsWithNotifications().catch(console.error)
