#!/usr/bin/env node

// Comprehensive test for notification action system
const baseUrl = 'http://localhost:8080'

async function testNotificationActionsWorkflow() {
  console.log('🧪 Testing complete notification action workflow...\n')

  try {
    // Step 1: Create a group invitation to generate a notification
    console.log('1️⃣ Creating a group invitation (which should generate a notification)...')
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
      console.log('   ✅ Group invitation created successfully')
    } else {
      const error = await inviteResponse.text()
      console.log(`   ⚠️ Invite response: ${inviteResponse.status} - ${error}`)
    }

    // Step 2: Fetch notifications to see what was created
    console.log('\n2️⃣ Fetching notifications to find invitation...')
    const notificationsResponse = await fetch(`${baseUrl}/notifications/fetch`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        start: 0,
        n_items: 10,
        type: 'all'
      })
    })

    if (notificationsResponse.ok) {
      const notificationData = await notificationsResponse.json()
      console.log(`   ✅ Found ${notificationData.notifications?.length || 0} notifications`)
      
      // Find a group invite notification
      const groupInviteNotification = notificationData.notifications?.find(n => 
        n.type === 'group_invite' && !n.seen
      )

      if (groupInviteNotification) {
        console.log(`   📧 Found group invite notification (ID: ${groupInviteNotification.id})`)
        console.log(`   Message: "${groupInviteNotification.message}"`)

        // Step 3: Test accepting the invitation via notification action
        console.log('\n3️⃣ Testing ACCEPT action via notification...')
        const acceptResponse = await fetch(`${baseUrl}/notifications/action`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            notification_id: groupInviteNotification.id,
            action: 'accept',
            group_id: 1  // Provide the group ID explicitly
          })
        })

        if (acceptResponse.ok) {
          const acceptData = await acceptResponse.json()
          console.log('   ✅ Successfully accepted invitation via notification!')
          console.log('   Response:', JSON.stringify(acceptData, null, 2))
        } else {
          const acceptError = await acceptResponse.text()
          console.log(`   ❌ Accept action failed: ${acceptResponse.status} - ${acceptError}`)
        }

      } else {
        console.log('   ❌ No unread group invite notifications found')
      }

    } else {
      const error = await notificationsResponse.text()
      console.log(`   ❌ Failed to fetch notifications: ${notificationsResponse.status} - ${error}`)
    }

    // Step 4: Test creating a follow request notification
    console.log('\n4️⃣ Testing follow request notification action...')
    const followResponse = await fetch(`${baseUrl}/profile/follow`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        target_id: 3
      })
    })

    if (followResponse.ok) {
      console.log('   ✅ Follow request created')

      // Fetch notifications again to find the follow request
      const followNotificationsResponse = await fetch(`${baseUrl}/notifications/fetch`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          start: 0,
          n_items: 10,
          type: 'all'
        })
      })

      if (followNotificationsResponse.ok) {
        const followNotificationData = await followNotificationsResponse.json()
        const followRequestNotification = followNotificationData.notifications?.find(n => 
          n.type === 'follow_request' && !n.seen
        )

        if (followRequestNotification) {
          console.log(`   📧 Found follow request notification (ID: ${followRequestNotification.id})`)
          console.log(`   Message: "${followRequestNotification.message}"`)

          // Test accepting the follow request
          console.log('\n5️⃣ Testing ACCEPT action for follow request...')
          const acceptFollowResponse = await fetch(`${baseUrl}/notifications/action`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              notification_id: followRequestNotification.id,
              action: 'accept',
              user_id: 1  // The user who sent the follow request
            })
          })

          if (acceptFollowResponse.ok) {
            const acceptFollowData = await acceptFollowResponse.json()
            console.log('   ✅ Successfully accepted follow request via notification!')
            console.log('   Response:', JSON.stringify(acceptFollowData, null, 2))
          } else {
            const acceptFollowError = await acceptFollowResponse.text()
            console.log(`   ❌ Accept follow action failed: ${acceptFollowResponse.status} - ${acceptFollowError}`)
          }
        } else {
          console.log('   ❌ No unread follow request notifications found')
        }
      }
    } else {
      const followError = await followResponse.text()
      console.log(`   ❌ Failed to create follow request: ${followResponse.status} - ${followError}`)
    }

    // Step 6: Test notification endpoints
    console.log('\n6️⃣ Testing notification management endpoints...')
    
    // Test mark as seen
    const markSeenResponse = await fetch(`${baseUrl}/notifications/mark-all-seen`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' }
    })

    if (markSeenResponse.ok) {
      console.log('   ✅ Mark all as seen endpoint works')
    } else {
      console.log(`   ❌ Mark as seen failed: ${markSeenResponse.status}`)
    }

    // Test unread count
    const unreadCountResponse = await fetch(`${baseUrl}/notifications/unread-count`)
    if (unreadCountResponse.ok) {
      const unreadData = await unreadCountResponse.json()
      console.log(`   ✅ Unread count: ${unreadData.unread_count || 0}`)
    } else {
      console.log(`   ❌ Unread count failed: ${unreadCountResponse.status}`)
    }

    console.log('\n🎉 Notification action workflow test completed!')

  } catch (error) {
    console.error('❌ Test failed with error:', error.message)
  }
}

// Run the test
testNotificationActionsWorkflow()
