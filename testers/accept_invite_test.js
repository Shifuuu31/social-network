// Test acceptGroupInvite and declineGroupInvite functions
const baseUrl = 'http://localhost:8080/groups/group'

console.log('🧪 Testing group invitation accept/decline functionality...\n')

async function testAcceptDeclineInvitations() {
  console.log('📧 Testing group invitation workflow...')
  
  try {
    // First, let's see the current state of user 3 (invited user)
    console.log('\n1️⃣ Checking current groups for user 3...')
    const initialResponse = await fetch(`${baseUrl}/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: "3",
        start: -1,
        n_items: 20,
        type: "user"
      })
    })
    const initialGroups = await initialResponse.json()
    console.log(`   User 3 currently has ${initialGroups?.length || 0} groups:`)
    initialGroups?.forEach(g => {
      console.log(`   - ${g.title} (${g.is_member})`)
    })
    
    // Look for a group where user 3 is invited
    const invitedGroup = initialGroups?.find(g => g.is_member === 'invited')
    
    if (invitedGroup) {
      console.log(`\n2️⃣ Found invitation to "${invitedGroup.title}" - Testing ACCEPT...`)
      
      // Test accepting the invitation
      const acceptResponse = await fetch(`${baseUrl}/accept-decline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: invitedGroup.id,
          user_id: 3,
          status: 'member',
          prev_status: 'invited'
        })
      })
      
      if (acceptResponse.ok) {
        const acceptData = await acceptResponse.json()
        console.log(`   ✅ Successfully accepted invitation!`)
        console.log(`   Response:`, acceptData)
        
        // Verify the change
        const verifyResponse = await fetch(`${baseUrl}/browse`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            user_id: "3",
            start: -1,
            n_items: 20,
            type: "user"
          })
        })
        const verifyGroups = await verifyResponse.json()
        const updatedGroup = verifyGroups?.find(g => g.id === invitedGroup.id)
        
        if (updatedGroup && updatedGroup.is_member === 'member') {
          console.log(`   ✅ Verification: User is now a member of "${updatedGroup.title}"`)
        } else {
          console.log(`   ❌ Verification failed: Status not updated properly`)
        }
      } else {
        const errorText = await acceptResponse.text()
        console.log(`   ❌ Failed to accept invitation: ${acceptResponse.status} ${errorText}`)
      }
    } else {
      console.log('\n2️⃣ No pending invitations found for user 3')
      
      // Create a test invitation first
      console.log('\n   Creating a test invitation...')
      const inviteResponse = await fetch(`${baseUrl}/invite`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: 4, // wesh a weldi group
          user_id: 3,
          status: 'invited',
          prev_status: 'none'
        })
      })
      
      if (inviteResponse.ok) {
        console.log(`   ✅ Test invitation created!`)
        
        // Now test accepting it
        console.log('\n   Testing accept invitation...')
        const acceptResponse = await fetch(`${baseUrl}/accept-decline`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            group_id: 4,
            user_id: 3,
            status: 'member',
            prev_status: 'invited'
          })
        })
        
        if (acceptResponse.ok) {
          const acceptData = await acceptResponse.json()
          console.log(`   ✅ Successfully accepted test invitation!`)
        } else {
          const errorText = await acceptResponse.text()
          console.log(`   ❌ Failed to accept test invitation: ${acceptResponse.status} ${errorText}`)
        }
      } else {
        const errorText = await inviteResponse.text()
        console.log(`   ❌ Failed to create test invitation: ${inviteResponse.status} ${errorText}`)
      }
    }
    
    // Test declining an invitation
    console.log('\n3️⃣ Testing DECLINE invitation...')
    
    // Create another test invitation to decline
    const declineTestResponse = await fetch(`${baseUrl}/invite`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 5, // ezjhbhezdjezhbazhjdbzdka group
        user_id: 3,
        status: 'invited',
        prev_status: 'none'
      })
    })
    
    if (declineTestResponse.ok) {
      console.log(`   📧 Test invitation created for decline test`)
      
      // Now decline it
      const declineResponse = await fetch(`${baseUrl}/accept-decline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: 5,
          user_id: 3,
          status: 'declined',
          prev_status: 'invited'
        })
      })
      
      if (declineResponse.ok) {
        console.log(`   ✅ Successfully declined invitation!`)
        
        // Verify the invitation was removed (should not appear in user's groups)
        const verifyResponse = await fetch(`${baseUrl}/browse`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            user_id: "3",
            start: -1,
            n_items: 20,
            type: "user"
          })
        })
        const verifyGroups = await verifyResponse.json()
        const declinedGroup = verifyGroups?.find(g => g.id === 5)
        
        if (!declinedGroup) {
          console.log(`   ✅ Verification: Declined invitation properly removed from user's groups`)
        } else {
          console.log(`   ❌ Verification failed: Declined group still appears in user's groups`)
        }
      } else {
        const errorText = await declineResponse.text()
        console.log(`   ❌ Failed to decline invitation: ${declineResponse.status} ${errorText}`)
      }
    }
    
    console.log('\n🎉 Accept/Decline invitation tests completed!')
    
  } catch (error) {
    console.error('❌ Test failed:', error.message)
  }
}

testAcceptDeclineInvitations()
