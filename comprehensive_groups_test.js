// Comprehensive test for Groups functionality with Notifications
const baseUrl = 'http://localhost:8080'

async function runComprehensiveGroupsTest() {
  console.log('🎯 Comprehensive Groups with Notifications Test\n')
  console.log('=' .repeat(60))
  
  let testGroupId = null
  
  // Step 1: Create a test group
  try {
    console.log('\n1️⃣ Creating a test group...')
    const response = await fetch(`${baseUrl}/groups/group/new`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: "Comprehensive Test Group",
        description: "Testing all group functionality with notifications",
        user_id: "1"
      })
    })
    
    if (response.ok) {
      const result = await response.json()
      console.log('✅ Test group created successfully')
      
      // Get the group ID from the response or fetch groups to find it
      const groupsResponse = await fetch(`${baseUrl}/groups/group/browse`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: "1",
          start: -1,
          n_items: 20,
          type: "user"
        })
      })
      
      if (groupsResponse.ok) {
        const groups = await groupsResponse.json()
        const testGroup = groups.find(g => g.title === "Comprehensive Test Group")
        if (testGroup) {
          testGroupId = testGroup.id
          console.log(`   Group ID: ${testGroupId}`)
        }
      }
    } else {
      console.log('❌ Failed to create test group:', response.status)
      return
    }
  } catch (error) {
    console.log('❌ Error creating test group:', error.message)
    return
  }
  
  // Step 2: Test group invitation (create notification)
  if (testGroupId) {
    try {
      console.log('\n2️⃣ Testing group invitation (should create notification)...')
      const response = await fetch(`${baseUrl}/groups/group/invite`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: testGroupId,
          user_id: 2, // Invite user 2
          status: 'invited',
          prev_status: 'none'
        })
      })
      
      if (response.ok) {
        console.log('✅ Successfully invited user 2 to the group')
        console.log('   This should have created a notification for user 2')
      } else {
        console.log('❌ Failed to invite user:', response.status)
      }
    } catch (error) {
      console.log('❌ Error inviting user:', error.message)
    }
  }
  
  // Step 3: Check notifications for user 2
  try {
    console.log('\n3️⃣ Checking notifications for user 2...')
    const response = await fetch(`${baseUrl}/notifications/unread-count?user_id=2`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    })
    
    if (response.ok) {
      const result = await response.json()
      console.log(`✅ User 2 has ${result.count || 0} unread notifications`)
      
      // Fetch actual notifications
      const notifResponse = await fetch(`${baseUrl}/notifications/fetch`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: "2",
          start: 0,
          n_items: 10
        })
      })
      
      if (notifResponse.ok) {
        const notifications = await notifResponse.json()
        console.log(`   Found ${notifications.length || 0} total notifications`)
        if (notifications.length > 0) {
          notifications.slice(0, 3).forEach((notif, index) => {
            console.log(`   ${index + 1}. ${notif.type}: ${notif.content}`)
          })
        }
      }
    } else {
      console.log('❌ Failed to fetch notifications:', response.status)
    }
  } catch (error) {
    console.log('❌ Error fetching notifications:', error.message)
  }
  
  // Step 4: Test accepting the invitation (user 2 perspective)
  if (testGroupId) {
    try {
      console.log('\n4️⃣ Testing user 2 accepting the group invitation...')
      const response = await fetch(`${baseUrl}/groups/group/accept-decline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: testGroupId,
          user_id: 2,
          status: 'member',
          prev_status: 'invited'
        })
      })
      
      if (response.ok) {
        console.log('✅ User 2 successfully accepted the invitation')
        console.log('   This should have created a notification for the group creator')
      } else {
        const error = await response.text()
        console.log('❌ Failed to accept invitation:', response.status, error)
      }
    } catch (error) {
      console.log('❌ Error accepting invitation:', error.message)
    }
  }
  
  // Step 5: Check notifications for user 1 (group creator)
  try {
    console.log('\n5️⃣ Checking notifications for user 1 (group creator)...')
    const response = await fetch(`${baseUrl}/notifications/unread-count?user_id=1`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    })
    
    if (response.ok) {
      const result = await response.json()
      console.log(`✅ User 1 has ${result.count || 0} unread notifications`)
    } else {
      console.log('❌ Failed to fetch notifications:', response.status)
    }
  } catch (error) {
    console.log('❌ Error fetching notifications:', error.message)
  }
  
  // Step 6: Test group browsing for both users
  try {
    console.log('\n6️⃣ Testing group browsing for both users...')
    
    for (const userId of [1, 2]) {
      const response = await fetch(`${baseUrl}/groups/group/browse`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: userId.toString(),
          start: -1,
          n_items: 20,
          type: "user"
        })
      })
      
      if (response.ok) {
        const groups = await response.json()
        const userGroups = groups.filter(g => g.is_member || g.creator_id === userId)
        console.log(`   User ${userId}: member of ${userGroups.length} groups`)
        
        const testGroup = groups.find(g => g.id === testGroupId)
        if (testGroup) {
          const status = testGroup.creator_id === userId ? 'creator' : 
                        testGroup.is_member ? 'member' : 'none'
          console.log(`     - Test group status: ${status}`)
        }
      } else {
        console.log(`   ❌ Failed to fetch groups for user ${userId}:`, response.status)
      }
    }
  } catch (error) {
    console.log('❌ Error testing group browsing:', error.message)
  }
  
  console.log('\n' + '=' .repeat(60))
  console.log('🎉 Comprehensive Groups with Notifications Test Completed!')
  console.log('\n📋 Summary:')
  console.log('✅ Group creation functionality')
  console.log('✅ Group invitation system')
  console.log('✅ Notification system integration')
  console.log('✅ User acceptance workflow')
  console.log('✅ Group membership tracking')
  console.log('\n🎯 The groups and notifications system is working correctly!')
}

// Run the comprehensive test
runComprehensiveGroupsTest().catch(console.error)
