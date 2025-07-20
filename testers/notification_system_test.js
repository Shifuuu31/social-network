// Test the notification system end-to-end
const baseUrl = 'http://localhost:8080'

console.log('🧪 Testing notification system...\n')

async function testNotificationSystem() {
  try {
    // Test 1: Get notifications (should be empty initially)
    console.log('1️⃣ Testing GET /notifications...')
    const notificationsResponse = await fetch(`${baseUrl}/notifications`, {
      method: 'GET',
      credentials: 'include'
    })
    
    if (notificationsResponse.ok) {
      const data = await notificationsResponse.json()
      console.log('✅ Notifications fetched successfully')
      console.log(`   Found ${data.notifications?.length || 0} notifications`)
      console.log('   Response:', JSON.stringify(data, null, 2))
    } else {
      console.log(`❌ Failed to fetch notifications: ${notificationsResponse.status}`)
      const error = await notificationsResponse.text()
      console.log(`   Error: ${error}`)
    }

    // Test 2: Get unread count
    console.log('\n2️⃣ Testing GET /notifications/unread-count...')
    const countResponse = await fetch(`${baseUrl}/notifications/unread-count`, {
      method: 'GET',
      credentials: 'include'
    })
    
    if (countResponse.ok) {
      const data = await countResponse.json()
      console.log('✅ Unread count fetched successfully')
      console.log(`   Unread count: ${data.unread_count}`)
    } else {
      console.log(`❌ Failed to fetch unread count: ${countResponse.status}`)
      const error = await countResponse.text()
      console.log(`   Error: ${error}`)
    }

    // Test 3: Create a group invitation to trigger notification
    console.log('\n3️⃣ Testing notification creation via group invite...')
    const inviteResponse = await fetch(`${baseUrl}/groups/group/invite`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        group_id: 1,
        user_id: 2,
        status: 'invited',
        prev_status: 'none'
      })
    })
    
    if (inviteResponse.ok) {
      console.log('✅ Group invite sent successfully (should create notification)')
    } else {
      console.log(`❌ Failed to send group invite: ${inviteResponse.status}`)
      const error = await inviteResponse.text()
      console.log(`   Error: ${error}`)
    }

    // Test 4: Check notifications again (should have new notification)
    console.log('\n4️⃣ Checking notifications after invite...')
    const newNotificationsResponse = await fetch(`${baseUrl}/notifications`, {
      method: 'GET',
      credentials: 'include'
    })
    
    if (newNotificationsResponse.ok) {
      const data = await newNotificationsResponse.json()
      console.log('✅ Notifications fetched successfully')
      console.log(`   Found ${data.notifications?.length || 0} notifications`)
      
      if (data.notifications && data.notifications.length > 0) {
        const latestNotification = data.notifications[0]
        console.log('   Latest notification:', JSON.stringify(latestNotification, null, 2))
        
        // Test 5: Mark notification as read
        if (latestNotification.id) {
          console.log('\n5️⃣ Testing mark notification as read...')
          const markReadResponse = await fetch(`${baseUrl}/notifications/mark-read`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({
              notification_ids: [latestNotification.id]
            })
          })
          
          if (markReadResponse.ok) {
            console.log('✅ Notification marked as read successfully')
          } else {
            console.log(`❌ Failed to mark notification as read: ${markReadResponse.status}`)
            const error = await markReadResponse.text()
            console.log(`   Error: ${error}`)
          }
        }
      }
    } else {
      console.log(`❌ Failed to fetch notifications: ${newNotificationsResponse.status}`)
    }

    console.log('\n🎉 Notification system tests completed!')
    
  } catch (error) {
    console.error('❌ Test failed:', error.message)
  }
}

testNotificationSystem()