// Comprehensive notification system test
const baseUrl = 'http://localhost:8080'

console.log('🔔 Testing Notification System...\n')

async function testNotificationSystem() {
  try {
    console.log('1️⃣ Testing Group Invitation Notification...')
    
    // Create a group invitation (this should trigger a notification)
    const inviteResponse = await fetch(`${baseUrl}/groups/group/invite`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 1,
        user_id: 2,
        status: 'invited',
        prev_status: 'none'
      })
    })
    
    if (inviteResponse.ok) {
      console.log('✅ Group invitation sent successfully')
    } else {
      console.log('❌ Failed to send group invitation:', await inviteResponse.text())
    }

    console.log('\n2️⃣ Testing Notification Fetch...')
    
    // Fetch notifications for user 2
    const notificationsResponse = await fetch(`${baseUrl}/notifications`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    })
    
    if (notificationsResponse.ok) {
      const notificationsData = await notificationsResponse.json()
      console.log('✅ Notifications fetched successfully')
      console.log('Notifications:', JSON.stringify(notificationsData, null, 2))
    } else {
      console.log('❌ Failed to fetch notifications:', await notificationsResponse.text())
    }

    console.log('\n3️⃣ Testing Unread Count...')
    
    // Get unread count
    const unreadResponse = await fetch(`${baseUrl}/notifications/unread-count`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    })
    
    if (unreadResponse.ok) {
      const unreadData = await unreadResponse.json()
      console.log('✅ Unread count fetched successfully')
      console.log('Unread count:', unreadData.unread_count)
    } else {
      console.log('❌ Failed to fetch unread count:', await unreadResponse.text())
    }

    console.log('\n4️⃣ Testing Join Request Notification...')
    
    // Create a join request (this should trigger a notification)
    const joinResponse = await fetch(`${baseUrl}/groups/group/request`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 2,
        user_id: 3,
        status: 'requested',
        prev_status: 'none'
      })
    })
    
    if (joinResponse.ok) {
      console.log('✅ Join request sent successfully')
    } else {
      console.log('❌ Failed to send join request:', await joinResponse.text())
    }

    console.log('\n5️⃣ Testing Event Creation Notification...')
    
    // Create an event (this should trigger notifications for all group members)
    const eventResponse = await fetch(`${baseUrl}/groups/group/event/new`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 1,
        title: 'Test Event',
        description: 'This is a test event for notification testing',
        event_time: new Date(Date.now() + 86400000).toISOString() // Tomorrow
      })
    })
    
    if (eventResponse.ok) {
      console.log('✅ Event created successfully')
    } else {
      console.log('❌ Failed to create event:', await eventResponse.text())
    }

    console.log('\n6️⃣ Testing Mark Notifications as Read...')
    
    // Mark all notifications as read
    const markReadResponse = await fetch(`${baseUrl}/notifications/mark-read`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        mark_all: true
      })
    })
    
    if (markReadResponse.ok) {
      console.log('✅ Notifications marked as read successfully')
    } else {
      console.log('❌ Failed to mark notifications as read:', await markReadResponse.text())
    }

    console.log('\n7️⃣ Testing Database Direct Query...')
    
    // This would require direct database access, which we can't do from frontend
    // But we can check if notifications were created by fetching again
    const finalNotificationsResponse = await fetch(`${baseUrl}/notifications`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    })
    
    if (finalNotificationsResponse.ok) {
      const finalData = await finalNotificationsResponse.json()
      console.log('✅ Final notification state:')
      console.log('Total notifications:', finalData.notifications?.length || 0)
      
      const unreadCount = finalData.notifications?.filter(n => !n.seen).length || 0
      console.log('Unread notifications:', unreadCount)
      
      // Show notification types
      const types = finalData.notifications?.map(n => n.type) || []
      console.log('Notification types found:', [...new Set(types)])
    }

    console.log('\n🎉 Notification System Test Completed!')
    
  } catch (error) {
    console.error('❌ Notification test failed:', error.message)
  }
}

// Run the test
testNotificationSystem()
