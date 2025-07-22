// Test the join group request functionality
const baseUrl = 'http://localhost:8080/groups/group'

console.log('🔄 Testing join group request workflow...')

async function testJoinGroupWorkflow() {
  const testUserId = 3
  const testGroupId = 1
  
  try {
    // Step 1: Check initial state
    console.log(`\n1️⃣ Checking initial groups for user ${testUserId}...`)
    const initialResponse = await fetch(`${baseUrl}/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: testUserId.toString(),
        start: -1,
        n_items: 20,
        type: "user"
      })
    })
    
    const initialGroups = await initialResponse.json()
    console.log(`   User ${testUserId} currently has ${initialGroups?.length || 0} groups:`)
    initialGroups?.forEach(g => {
      console.log(`   - Group ${g.id}: "${g.title}" (${g.is_member || 'none'})`)
    })
    
    const existingGroup = initialGroups?.find(g => g.id === testGroupId)
    if (existingGroup) {
      console.log(`   ⚠️ User ${testUserId} already has a relationship with group ${testGroupId}: ${existingGroup.is_member}`)
    }
    
    // Step 2: Send join request
    console.log(`\n2️⃣ User ${testUserId} requesting to join group ${testGroupId}...`)
    const requestResponse = await fetch(`${baseUrl}/request`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: testUserId,
        group_id: testGroupId,
        status: 'requested',
        prev_status: 'none'
      })
    })
    
    if (requestResponse.ok) {
      const requestData = await requestResponse.json()
      console.log(`   ✅ Join request successful! Status: ${requestData}`)
      
      // Step 3: Verify the change
      console.log(`\n3️⃣ Verifying user ${testUserId} groups after request...`)
      const verifyResponse = await fetch(`${baseUrl}/browse`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: testUserId.toString(),
          start: -1,
          n_items: 20,
          type: "user"
        })
      })
      
      if (verifyResponse.ok) {
        const updatedGroups = await verifyResponse.json()
        console.log(`   User ${testUserId} groups after request (${updatedGroups?.length || 0}):`)
        updatedGroups?.forEach(g => {
          console.log(`   - Group ${g.id}: "${g.title}" (${g.is_member || 'none'})`)
          if (g.id === testGroupId) {
            console.log(`     🎯 TARGET GROUP! Status changed to: ${g.is_member}`)
          }
        })
        
        const targetGroup = updatedGroups?.find(g => g.id === testGroupId)
        if (targetGroup && targetGroup.is_member === 'requested') {
          console.log(`   ✅ SUCCESS! User ${testUserId} now has 'requested' status for group ${testGroupId}`)
        } else {
          console.log(`   ❌ FAILED! Expected 'requested' status but got: ${targetGroup?.is_member || 'none'}`)
        }
      } else {
        console.log(`   ❌ Failed to verify groups: ${verifyResponse.status}`)
      }
    } else {
      const error = await requestResponse.text()
      console.log(`   ❌ Join request failed: ${requestResponse.status} - ${error}`)
    }
    
  } catch (error) {
    console.error('❌ Error in join group workflow test:', error.message)
  }
}

// Run the test
testJoinGroupWorkflow().catch(console.error)
